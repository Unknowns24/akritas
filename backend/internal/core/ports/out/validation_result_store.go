package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type ValidationResultStore interface {
	Create(ctx context.Context, value *domain.ValidationResult) error
	ListByRemediation(ctx context.Context, remediationID uuid.UUID) ([]domain.ValidationResult, error)
}
