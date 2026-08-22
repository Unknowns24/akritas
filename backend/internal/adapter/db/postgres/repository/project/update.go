package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Update(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}
