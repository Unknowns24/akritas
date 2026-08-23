package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type OperationUseCase interface {
	GetOperation(context.Context, uuid.UUID) (*domain.Operation, error)
}
