package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type verifyAdministratorSetupUseCase struct {
	rateLimiter        out.RateLimiter
	pendingEnrollments out.PendingEnrollmentRepository
	credentialStore    out.CredentialStore
	totpVerifier       out.TOTPVerifier
	administrators     out.AdministratorRepository
	sessionTokens      out.SessionTokenGenerator
	sessions           out.AdministratorSessionRepository
	transactor         out.Transactor
	now                func() time.Time
	sessionIdleTTL     time.Duration
	sessionAbsoluteTTL time.Duration
}

func NewVerifyAdministratorSetupUseCase(
	rateLimiter out.RateLimiter,
	pendingEnrollments out.PendingEnrollmentRepository,
	credentialStore out.CredentialStore,
	totpVerifier out.TOTPVerifier,
	administrators out.AdministratorRepository,
	sessionTokens out.SessionTokenGenerator,
	sessions out.AdministratorSessionRepository,
	transactor out.Transactor,
	now func() time.Time,
	sessionIdleTTL time.Duration,
	sessionAbsoluteTTL time.Duration,
) in.VerifyAdministratorSetupUseCase {
	return &verifyAdministratorSetupUseCase{
		rateLimiter:        rateLimiter,
		pendingEnrollments: pendingEnrollments,
		credentialStore:    credentialStore,
		totpVerifier:       totpVerifier,
		administrators:     administrators,
		sessionTokens:      sessionTokens,
		sessions:           sessions,
		transactor:         transactor,
		now:                now,
		sessionIdleTTL:     sessionIdleTTL,
		sessionAbsoluteTTL: sessionAbsoluteTTL,
	}
}

func (uc *verifyAdministratorSetupUseCase) Execute(ctx context.Context, input in.VerifyAdministratorSetupInput) (in.VerifyAdministratorSetupOutput, error) {
	allowed, err := uc.rateLimiter.Allow(ctx, input.RateLimitKey)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}
	if !allowed {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrAuthenticationRateLimited
	}

	enrollmentID, err := uuid.Parse(input.EnrollmentID)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrInvalidTotpEnrollmentVerification
	}

	enrollment, err := uc.pendingEnrollments.FindByID(ctx, enrollmentID)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}
	if enrollment == nil {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrInvalidTotpEnrollmentVerification
	}

	now := uc.now().UTC()
	if enrollment.Enrollment.IsExpired(now) {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrInvalidTotpEnrollmentVerification
	}

	secret, err := uc.credentialStore.Get(ctx, out.CredentialOwnerPendingEnrollment, enrollmentID, out.SecretKindAdministratorTOTP)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrInvalidTotpEnrollmentVerification
	}
	defer clear(secret)

	valid, _, err := uc.totpVerifier.Verify(string(secret), input.TOTPCode, now)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}
	if !valid {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrInvalidTotpEnrollmentVerification
	}

	exists, err := uc.administrators.ExistsActive(ctx)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}
	if exists {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrAdministratorAlreadyExists
	}

	administrator, err := domain.NewAdministrator(uuid.New(), enrollment.Enrollment.Email, enrollment.Enrollment.DisplayName, now)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}

	sessionToken, tokenHash, err := uc.sessionTokens.Generate()
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}

	session, err := domain.NewAdministratorSession(uuid.New(), administrator.ID, now, now.Add(uc.sessionIdleTTL), now.Add(uc.sessionAbsoluteTTL))
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}

	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		current, findErr := uc.pendingEnrollments.FindByID(txCtx, enrollmentID)
		if findErr != nil {
			return findErr
		}
		if current == nil || current.Enrollment.IsExpired(now) {
			return domain.ErrInvalidTotpEnrollmentVerification
		}
		if err := uc.administrators.Create(txCtx, administrator, current.PasswordHash); err != nil {
			return err
		}
		if err := uc.sessions.Save(txCtx, session, tokenHash); err != nil {
			return err
		}
		if err := uc.credentialStore.MoveOwner(txCtx, out.CredentialOwnerPendingEnrollment, enrollmentID, out.CredentialOwnerAdministrator, administrator.ID); err != nil {
			return err
		}
		return uc.pendingEnrollments.Delete(txCtx, enrollmentID)
	})
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}

	return in.VerifyAdministratorSetupOutput{
		Administrator:     *administrator,
		SessionToken:      sessionToken,
		AuthenticatedAt:   session.AuthenticatedAt,
		IdleExpiresAt:     session.IdleExpiresAt,
		AbsoluteExpiresAt: session.AbsoluteExpiresAt,
	}, nil
}
