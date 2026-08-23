package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Investigation, error) {
	var value domain.Investigation
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("investigations").Where("id = ?", id).First(&value).Error; err != nil {
		return nil, mapError(err)
	}
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
