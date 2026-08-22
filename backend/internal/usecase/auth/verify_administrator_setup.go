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
	clock              out.Clock
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
	clock out.Clock,
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
		clock:              clock,
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
		return in.VerifyAdministratorSetupOutput{}, ErrSetupRateLimited
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

	now := uc.clock.Now()
	if enrollment.IsExpired(now) {
		return in.VerifyAdministratorSetupOutput{}, domain.ErrInvalidTotpEnrollmentVerification
	}

	secret, err := uc.credentialStore.Decrypt(ctx, enrollment.EncryptedTOTPSecret)
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}

	valid, err := uc.totpVerifier.Verify(string(secret), input.TOTPCode, now)
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

	administrator, err := domain.NewAdministrator(uuid.New(), enrollment.Email, enrollment.DisplayName, now)
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

	err = uc.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := uc.administrators.Create(ctx, administrator, enrollment.PasswordHash, enrollment.EncryptedTOTPSecret); err != nil {
			return err
		}
		return uc.sessions.Save(ctx, session, tokenHash)
	})
	if err != nil {
		return in.VerifyAdministratorSetupOutput{}, err
	}

	if err := uc.pendingEnrollments.Delete(ctx, enrollment.ID); err != nil {
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
