package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// OperationStore is generic async-command persistence, shared by every
// resource type that queues work through Operation (investigation today;
// remediation and pull_request are expected to reuse it later).
type OperationStore interface {
	Create(context.Context, *domain.Operation) error
	Update(context.Context, *domain.Operation) error
	FindByID(context.Context, uuid.UUID) (*domain.Operation, error)
	FindByIdempotencyKey(context.Context, string) (*domain.Operation, error)
}
