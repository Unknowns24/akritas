package validationresult

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Create(ctx context.Context, value *domain.ValidationResult) error {
	record := fromDomain(value)
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("validation_results").Create(&record).Error; err != nil {
		return mapError(err)
	}
	return nil
}
