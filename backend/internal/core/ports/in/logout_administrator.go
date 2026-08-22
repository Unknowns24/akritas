package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// LogoutAdministratorUseCase takes a session already resolved by
// AuthenticateSessionUseCase and revokes it.
type LogoutAdministratorUseCase interface {
	Execute(ctx context.Context, session domain.AdministratorSession) error
}
