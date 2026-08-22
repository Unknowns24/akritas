package auth

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type logoutAdministratorUseCase struct {
	sessions out.AdministratorSessionRepository
	clock    out.Clock
}

func NewLogoutAdministratorUseCase(sessions out.AdministratorSessionRepository, clock out.Clock) in.LogoutAdministratorUseCase {
	return &logoutAdministratorUseCase{sessions: sessions, clock: clock}
}

func (uc *logoutAdministratorUseCase) Execute(ctx context.Context, session domain.AdministratorSession) error {
	now := uc.clock.Now()
	if err := session.Revoke(now); err != nil {
		return err
	}
	return uc.sessions.Revoke(ctx, session.ID, *session.RevokedAt)
}
