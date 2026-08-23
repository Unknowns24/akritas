package operation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByResource(ctx context.Context, resourceType domain.OperationResourceType, resourceID uuid.UUID) (*domain.Operation, error) {
	var value domain.Operation
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("operations").
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC, id DESC").Take(&value).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
