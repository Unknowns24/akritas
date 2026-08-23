package monitoring

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) ListProjectsForMonitoring(ctx context.Context) ([]domain.Project, error) {
	var projects []domain.Project
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").
		Joins("LEFT JOIN monitoring_checkpoints ON monitoring_checkpoints.project_id = projects.id AND monitoring_checkpoints.is_current").
		Where("projects.monitoring_enabled = ? OR monitoring_checkpoints.next_finalize_at IS NOT NULL", true).
		Order("projects.id ASC").Find(&projects).Error
	if err != nil {
		return nil, mapError(err)
	}
	return projects, nil
}
