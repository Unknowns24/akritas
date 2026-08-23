package auth

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type verifyAdministratorRecoveryUseCase struct {
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

func NewVerifyAdministratorRecoveryUseCase(rateLimiter out.RateLimiter, pendingEnrollments out.PendingEnrollmentRepository, credentialStore out.CredentialStore, totpVerifier out.TOTPVerifier, administrators out.AdministratorRepository, sessionTokens out.SessionTokenGenerator, sessions out.AdministratorSessionRepository, transactor out.Transactor, now func() time.Time, sessionIdleTTL, sessionAbsoluteTTL time.Duration) in.VerifyAdministratorRecoveryUseCase {
	return &verifyAdministratorRecoveryUseCase{rateLimiter: rateLimiter, pendingEnrollments: pendingEnrollments, credentialStore: credentialStore, totpVerifier: totpVerifier, administrators: administrators, sessionTokens: sessionTokens, sessions: sessions, transactor: transactor, now: now, sessionIdleTTL: sessionIdleTTL, sessionAbsoluteTTL: sessionAbsoluteTTL}
}

func (uc *verifyAdministratorRecoveryUseCase) Execute(ctx context.Context, input in.VerifyAdministratorRecoveryInput) (in.VerifyAdministratorRecoveryOutput, error) {
	for _, key := range []string{"ip:" + input.RateLimitKey, "enrollment:" + input.EnrollmentID} {
		allowed, err := uc.rateLimiter.Allow(ctx, key)
		if err != nil {
			return in.VerifyAdministratorRecoveryOutput{}, err
		}
		if !allowed {
			return in.VerifyAdministratorRecoveryOutput{}, domain.ErrAuthenticationRateLimited
		}
	}
	enrollmentID, err := uuid.Parse(input.EnrollmentID)
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, domain.ErrInvalidCredentials
	}
	now := uc.now().UTC()
	enrollment, err := uc.pendingEnrollments.FindByID(ctx, enrollmentID)
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, err
	}
	secret := []byte(dummyTOTPSecret)
	credentialOK := false
	if enrollment != nil && !enrollment.Enrollment.IsExpired(now) {
		storedSecret, getErr := uc.credentialStore.Get(ctx, out.CredentialOwnerPendingEnrollment, enrollmentID, out.SecretKindAdministratorTOTP)
		if getErr == nil {
			secret = storedSecret
			credentialOK = true
		}
	}
	defer clear(secret)
	valid, period, err := uc.totpVerifier.Verify(string(secret), input.TOTPCode, now)
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, err
	}
	if enrollment == nil || enrollment.Enrollment.IsExpired(now) || !credentialOK || !valid {
		return in.VerifyAdministratorRecoveryOutput{}, domain.ErrInvalidCredentials
	}
	current, err := uc.administrators.FindByEmail(ctx, enrollment.Enrollment.Email)
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, err
	}
	if current == nil {
		return in.VerifyAdministratorRecoveryOutput{}, domain.ErrInvalidCredentials
	}
	sessionToken, tokenHash, err := uc.sessionTokens.Generate()
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, err
	}
	session, err := domain.NewAdministratorSession(uuid.New(), current.Administrator.ID, now, now.Add(uc.sessionIdleTTL), now.Add(uc.sessionAbsoluteTTL))
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, err
	}

	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		consumed, consumeErr := uc.pendingEnrollments.Consume(txCtx, enrollmentID)
		if consumeErr != nil {
			return consumeErr
		}
		if consumed == nil || consumed.Enrollment.IsExpired(now) || consumed.Enrollment.Email != current.Administrator.Email || consumed.PasswordHash != enrollment.PasswordHash {
			return domain.ErrInvalidCredentials
		}
		rotated, rotateErr := uc.administrators.RotateCredentials(txCtx, current.Administrator.ID, current.PasswordHash, consumed.PasswordHash, period, now)
		if rotateErr != nil {
			return rotateErr
		}
		if !rotated {
			return domain.ErrInvalidCredentials
		}
		if moveErr := uc.credentialStore.MoveOwner(txCtx, out.CredentialOwnerPendingEnrollment, enrollmentID, out.CredentialOwnerAdministrator, current.Administrator.ID); moveErr != nil {
			return moveErr
		}
		if revokeErr := uc.sessions.RevokeAll(txCtx, current.Administrator.ID, now); revokeErr != nil {
			return revokeErr
		}
		return uc.sessions.Save(txCtx, session, tokenHash)
	})
	if err != nil {
		return in.VerifyAdministratorRecoveryOutput{}, err
	}
	administrator := current.Administrator
	administrator.LastAcceptedTOTPPeriod = period
	administrator.UpdatedAt = now
	return in.VerifyAdministratorRecoveryOutput{Administrator: administrator, SessionToken: sessionToken, AuthenticatedAt: session.AuthenticatedAt, IdleExpiresAt: session.IdleExpiresAt, AbsoluteExpiresAt: session.AbsoluteExpiresAt}, nil
}
