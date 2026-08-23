package monitoring

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) UpdateProjectObservation(ctx context.Context, projectID uuid.UUID, status domain.MonitoringStatus, updatedAt time.Time, observedAt *time.Time) error {
	values := map[string]any{"monitoring_status": status, "updated_at": updatedAt.UTC()}
	if observedAt != nil {
		values["last_observed_at"] = observedAt.UTC()
	}
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("id = ?", projectID).Updates(values)
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrProjectNotFound
	}
	return nil
}
