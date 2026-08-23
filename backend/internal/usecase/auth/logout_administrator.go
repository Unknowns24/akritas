package auth

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type logoutAdministratorUseCase struct {
	sessions out.AdministratorSessionRepository
	now      func() time.Time
}

func NewLogoutAdministratorUseCase(sessions out.AdministratorSessionRepository, now func() time.Time) in.LogoutAdministratorUseCase {
	return &logoutAdministratorUseCase{sessions: sessions, now: now}
}

func (uc *logoutAdministratorUseCase) Execute(ctx context.Context, session domain.AdministratorSession) error {
	now := uc.now().UTC()
	if err := session.Revoke(now); err != nil {
		return err
	}
	return uc.sessions.Revoke(ctx, session.ID, *session.RevokedAt)
}
