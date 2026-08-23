package evidence

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Create(ctx context.Context, value *domain.Evidence) error {
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("evidence").Create(value).Error; err != nil {
		return mapError(err)
	}
	return nil
}
