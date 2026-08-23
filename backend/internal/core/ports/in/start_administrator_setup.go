package in

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StartAdministratorSetupInput struct {
	Email          string
	DisplayName    string
	Password       string
	BootstrapToken string
	RateLimitKey   string
}

type StartAdministratorSetupOutput struct {
	EnrollmentID   uuid.UUID
	OtpauthURI     string
	ManualEntryKey string
	ExpiresAt      time.Time
}

type StartAdministratorSetupUseCase interface {
	Execute(ctx context.Context, input StartAdministratorSetupInput) (StartAdministratorSetupOutput, error)
}
