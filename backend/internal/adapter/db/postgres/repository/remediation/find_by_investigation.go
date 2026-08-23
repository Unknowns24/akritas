package remediation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByInvestigation(ctx context.Context, investigationID uuid.UUID) (*domain.Remediation, error) {
	if investigationID == uuid.Nil {
		return nil, domain.ErrInvalidRemediation
	}
	var record remediationRecord
	if err := txcontext.From(ctx, r.db).WithContext(ctx).
		Table("remediations").
		Where("investigation_id = ?", investigationID).
		First(&record).Error; err != nil {
		return nil, mapError(err)
	}
	value := record.toDomain()
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return value, nil
}
