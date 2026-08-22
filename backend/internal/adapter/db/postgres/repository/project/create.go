package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Create(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}
