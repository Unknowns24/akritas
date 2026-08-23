package in

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type GetCurrentSessionOutput struct {
	Administrator     domain.Administrator
	AuthenticatedAt   time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

// GetCurrentSessionUseCase takes a session already resolved (and
// idle-extended) by AuthenticateSessionUseCase and builds the safe
// projection to return.
type GetCurrentSessionUseCase interface {
	Execute(ctx context.Context, session domain.AdministratorSession) (GetCurrentSessionOutput, error)
}
