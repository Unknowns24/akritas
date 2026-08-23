package auth

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type loginAdministratorUseCase struct {
	rateLimiter        out.RateLimiter
	administrators     out.AdministratorRepository
	passwordHasher     out.PasswordHasher
	credentialStore    out.CredentialStore
	totpVerifier       out.TOTPVerifier
	sessionTokens      out.SessionTokenGenerator
	sessions           out.AdministratorSessionRepository
	transactor         out.Transactor
	now                func() time.Time
	sessionIdleTTL     time.Duration
	sessionAbsoluteTTL time.Duration
}

func NewLoginAdministratorUseCase(
	rateLimiter out.RateLimiter,
	administrators out.AdministratorRepository,
	passwordHasher out.PasswordHasher,
	credentialStore out.CredentialStore,
	totpVerifier out.TOTPVerifier,
	sessionTokens out.SessionTokenGenerator,
	sessions out.AdministratorSessionRepository,
	transactor out.Transactor,
	now func() time.Time,
	sessionIdleTTL time.Duration,
	sessionAbsoluteTTL time.Duration,
) in.LoginAdministratorUseCase {
	return &loginAdministratorUseCase{
		rateLimiter:        rateLimiter,
		administrators:     administrators,
		passwordHasher:     passwordHasher,
		credentialStore:    credentialStore,
		totpVerifier:       totpVerifier,
		sessionTokens:      sessionTokens,
		sessions:           sessions,
		transactor:         transactor,
		now:                now,
		sessionIdleTTL:     sessionIdleTTL,
		sessionAbsoluteTTL: sessionAbsoluteTTL,
	}
}

func (uc *loginAdministratorUseCase) Execute(ctx context.Context, input in.LoginAdministratorInput) (in.LoginAdministratorOutput, error) {
	ipAllowed, err := uc.rateLimiter.Allow(ctx, "ip:"+input.RateLimitKey)
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}
	if !ipAllowed {
		return in.LoginAdministratorOutput{}, domain.ErrAuthenticationRateLimited
	}

	accountAllowed, err := uc.rateLimiter.Allow(ctx, "account:"+strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}
	if !accountAllowed {
		return in.LoginAdministratorOutput{}, domain.ErrAuthenticationRateLimited
	}

	creds, err := uc.administrators.FindByEmail(ctx, input.Email)
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}
	if creds == nil {
		return in.LoginAdministratorOutput{}, domain.ErrInvalidCredentials
	}

	passwordOK, err := uc.passwordHasher.Verify(input.Password, creds.PasswordHash)
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}
	if !passwordOK {
		return in.LoginAdministratorOutput{}, domain.ErrInvalidCredentials
	}

	secret, err := uc.credentialStore.Get(ctx, out.CredentialOwnerAdministrator, creds.Administrator.ID, out.SecretKindAdministratorTOTP)
	if err != nil {
		return in.LoginAdministratorOutput{}, domain.ErrInvalidCredentials
	}
	defer clear(secret)

	now := uc.now().UTC()
	valid, period, err := uc.totpVerifier.Verify(string(secret), input.TOTPCode, now)
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}
	if !valid {
		return in.LoginAdministratorOutput{}, domain.ErrInvalidCredentials
	}
	if period <= creds.Administrator.LastAcceptedTOTPPeriod {
		return in.LoginAdministratorOutput{}, domain.ErrInvalidCredentials
	}

	sessionToken, tokenHash, err := uc.sessionTokens.Generate()
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}

	session, err := domain.NewAdministratorSession(uuid.New(), creds.Administrator.ID, now, now.Add(uc.sessionIdleTTL), now.Add(uc.sessionAbsoluteTTL))
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}

	err = uc.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		consumed, err := uc.administrators.ConsumeTOTPPeriod(ctx, creds.Administrator.ID, period)
		if err != nil {
			return err
		}
		if !consumed {
			return domain.ErrInvalidCredentials
		}
		return uc.sessions.Save(ctx, session, tokenHash)
	})
	if err != nil {
		return in.LoginAdministratorOutput{}, err
	}

	return in.LoginAdministratorOutput{
		Administrator:     creds.Administrator,
		SessionToken:      sessionToken,
		AuthenticatedAt:   session.AuthenticatedAt,
		IdleExpiresAt:     session.IdleExpiresAt,
		AbsoluteExpiresAt: session.AbsoluteExpiresAt,
	}, nil
}
