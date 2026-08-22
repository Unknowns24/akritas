package auth

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type authenticateSessionUseCase struct {
	sessions       out.AdministratorSessionRepository
	sessionTokens  out.SessionTokenGenerator
	clock          out.Clock
	sessionIdleTTL time.Duration
}

func NewAuthenticateSessionUseCase(
	sessions out.AdministratorSessionRepository,
	sessionTokens out.SessionTokenGenerator,
	clock out.Clock,
	sessionIdleTTL time.Duration,
) in.AuthenticateSessionUseCase {
	return &authenticateSessionUseCase{
		sessions:       sessions,
		sessionTokens:  sessionTokens,
		clock:          clock,
		sessionIdleTTL: sessionIdleTTL,
	}
}

func (uc *authenticateSessionUseCase) Execute(ctx context.Context, sessionToken string) (domain.AdministratorSession, error) {
	if sessionToken == "" {
		return domain.AdministratorSession{}, domain.ErrInactiveAdministratorSession
	}

	tokenHash := uc.sessionTokens.Hash(sessionToken)
	session, err := uc.sessions.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return domain.AdministratorSession{}, err
	}
	if session == nil {
		return domain.AdministratorSession{}, domain.ErrInactiveAdministratorSession
	}

	now := uc.clock.Now()
	if !session.IsActive(now) {
		return domain.AdministratorSession{}, domain.ErrInactiveAdministratorSession
	}

	if err := session.ExtendIdle(now, uc.sessionIdleTTL); err != nil {
		return domain.AdministratorSession{}, err
	}
	if err := uc.sessions.UpdateIdleExpiry(ctx, session.ID, session.IdleExpiresAt); err != nil {
		return domain.AdministratorSession{}, err
	}

	return *session, nil
}
