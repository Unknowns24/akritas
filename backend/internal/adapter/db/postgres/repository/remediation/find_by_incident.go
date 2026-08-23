package remediation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByIncident(ctx context.Context, incidentID uuid.UUID) (*domain.Remediation, error) {
	if incidentID == uuid.Nil {
		return nil, domain.ErrInvalidRemediation
	}
	var record remediationRecord
	if err := txcontext.From(ctx, r.db).WithContext(ctx).
		Table("remediations").
		Where("incident_id = ?", incidentID).
		Order("created_at DESC, id DESC").
		First(&record).Error; err != nil {
		return nil, mapError(err)
	}
	value := record.toDomain()
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return value, nil
}
