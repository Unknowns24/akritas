package in

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StartAdministratorRecoveryInput struct {
	Email          string
	NewPassword    string
	BootstrapToken string
	RateLimitKey   string
}

type StartAdministratorRecoveryOutput struct {
	EnrollmentID   uuid.UUID
	OtpauthURI     string
	ManualEntryKey string
	ExpiresAt      time.Time
}

type StartAdministratorRecoveryUseCase interface {
	Execute(context.Context, StartAdministratorRecoveryInput) (StartAdministratorRecoveryOutput, error)
}
