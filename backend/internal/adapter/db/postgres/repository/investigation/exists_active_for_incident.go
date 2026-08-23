package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) ExistsActiveForIncident(ctx context.Context, incidentID uuid.UUID) (bool, error) {
	var count int64
	active := []string{string(domain.InvestigationStatusPending), string(domain.InvestigationStatusRunning)}
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("investigations").
		Where("incident_id = ? AND status IN ?", incidentID, active).Count(&count).Error; err != nil {
		return false, mapError(err)
	}
	return count > 0, nil
}
