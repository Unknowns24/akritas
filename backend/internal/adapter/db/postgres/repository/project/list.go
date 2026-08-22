package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

func (r *Repository) List(ctx context.Context, query paging.ListQuery) ([]domain.Project, int64, error) {
	db := r.db.WithContext(ctx).Model(&domain.Project{})
	if query.NameLike != "" {
		db = db.Where("LOWER(name) LIKE ?", "%"+query.NameLike+"%")
	}
	if len(query.MonitoringStatusIn) > 0 {
		statuses := make([]string, 0, len(query.MonitoringStatusIn))
		for _, status := range query.MonitoringStatusIn {
			statuses = append(statuses, string(status))
		}
		db = db.Where("monitoring_status IN ?", statuses)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := "created_at DESC, id DESC"
	if query.Sort == "name:asc" {
		order = "name ASC, id ASC"
	}
	var projects []domain.Project
	if err := db.Order(order).Offset(query.Offset).Limit(query.Limit).Find(&projects).Error; err != nil {
		return nil, 0, err
	}
	for i := range projects {
		if err := projects[i].Validate(); err != nil {
			return nil, 0, err
		}
	}
	return projects, total, nil
}
