package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	var row incidentRow
	db := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").
		Select("incidents.*, (SELECT name FROM projects WHERE projects.id = incidents.project_id) AS project_name").
		Where("incidents.id = ?", id).Take(&row)
	if db.Error != nil {
		return nil, mapError(db.Error)
	}
	value := row.domain()
	return &value, nil
}

type incidentRow struct {
	domain.Incident `gorm:"embedded"`
	ProjectName     string `gorm:"column:project_name"`
}

func (r incidentRow) domain() domain.Incident {
	value := r.Incident
	value.Project = &domain.ProjectReference{ID: value.ProjectID, Name: r.ProjectName}
	return value
}
