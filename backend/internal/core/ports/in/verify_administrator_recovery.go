package in

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type VerifyAdministratorRecoveryInput struct {
	EnrollmentID string
	TOTPCode     string
	RateLimitKey string
}

type VerifyAdministratorRecoveryOutput struct {
	Administrator     domain.Administrator
	SessionToken      string
	AuthenticatedAt   time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

type VerifyAdministratorRecoveryUseCase interface {
	Execute(context.Context, VerifyAdministratorRecoveryInput) (VerifyAdministratorRecoveryOutput, error)
}
