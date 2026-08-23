package project

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Update(ctx context.Context, value *domain.Project, expected time.Time) error {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("id = ? AND updated_at = ?", value.ID, expected).Select("*").Omit("id", "created_at").Updates(value)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrProjectConcurrentModification
	}
	return nil
}
