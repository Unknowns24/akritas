package in

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type VerifyAdministratorSetupInput struct {
	EnrollmentID string
	TOTPCode     string
	RateLimitKey string
}

type VerifyAdministratorSetupOutput struct {
	Administrator     domain.Administrator
	SessionToken      string
	AuthenticatedAt   time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

type VerifyAdministratorSetupUseCase interface {
	Execute(ctx context.Context, input VerifyAdministratorSetupInput) (VerifyAdministratorSetupOutput, error)
}
