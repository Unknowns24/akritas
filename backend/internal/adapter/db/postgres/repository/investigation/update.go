package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Update(ctx context.Context, value *domain.Investigation) error {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("investigations").Where("id = ?", value.ID).
		Select("*").Omit("id", "created_at", "incident_id").Updates(value)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrInvestigationNotFound
	}
	return nil
}
