package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) ListOpen(ctx context.Context) ([]domain.Investigation, error) {
	var values []domain.Investigation
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("investigations").
		Where("status IN ?", []domain.InvestigationStatus{domain.InvestigationStatusPending, domain.InvestigationStatusRunning}).
		Order("created_at ASC, id ASC").Find(&values).Error
	if err != nil {
		return nil, mapError(err)
	}
	return values, nil
}
