package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// AuthenticateSessionUseCase resolves the raw cookie value into an active,
// idle-extended session. Used by the REST auth middleware; sessionToken=""
// (no cookie) resolves to domain.ErrInactiveAdministratorSession without
// touching the database.
type AuthenticateSessionUseCase interface {
	Execute(ctx context.Context, sessionToken string) (domain.AdministratorSession, error)
}
