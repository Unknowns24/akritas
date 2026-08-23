package monitoring

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (r *Repository) LockProject(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Take(&project).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &project, nil
}
