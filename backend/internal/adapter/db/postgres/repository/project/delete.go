package project

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("id = ? AND monitoring_enabled = false", id).Delete(&domain.Project{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
			return domain.ErrProjectHasDependencies
		}
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrProjectMustBeDisabled
	}
	return nil
}
