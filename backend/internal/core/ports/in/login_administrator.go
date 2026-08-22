package in

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type LoginAdministratorInput struct {
	Email        string
	Password     string
	TOTPCode     string
	RateLimitKey string
}

type LoginAdministratorOutput struct {
	Administrator     domain.Administrator
	SessionToken      string
	AuthenticatedAt   time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

type LoginAdministratorUseCase interface {
	Execute(ctx context.Context, input LoginAdministratorInput) (LoginAdministratorOutput, error)
}
