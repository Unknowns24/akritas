package validationresult

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) ListByRemediation(ctx context.Context, remediationID uuid.UUID) ([]domain.ValidationResult, error) {
	var records []validationResultRecord
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("validation_results").
		Where("remediation_id = ?", remediationID).
		Order("created_at ASC, id ASC").
		Find(&records).Error; err != nil {
		return nil, mapError(err)
	}
	results := make([]domain.ValidationResult, 0, len(records))
	for _, record := range records {
		results = append(results, record.toDomain())
	}
	return results, nil
}
