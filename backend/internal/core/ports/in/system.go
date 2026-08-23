package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type SystemUseCase interface {
	GetStatus(context.Context) (domain.SystemStatus, error)
	RunDiagnostics(context.Context, uuid.UUID) (*domain.Operation, error)
}
