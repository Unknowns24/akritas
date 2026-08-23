package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) ListByIncident(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.Investigation], error) {
	db := txcontext.From(ctx, r.db).WithContext(ctx).Table("investigations").Model(&domain.Investigation{}).Where("incident_id = ?", incidentID)
	query, err := ukerpagination.Apply(db, params)
	if err != nil {
		return paging.Slice[domain.Investigation]{}, domain.ErrInvalidInvestigation.Wrap(err)
	}
	var items []domain.Investigation
	if err := query.Find(&items).Error; err != nil {
		return paging.Slice[domain.Investigation]{}, mapError(err)
	}
	for i := range items {
		if err := items[i].Validate(); err != nil {
			return paging.Slice[domain.Investigation]{}, mapError(err)
		}
	}
	countBase := txcontext.From(ctx, r.db).WithContext(ctx).Table("investigations").Model(&domain.Investigation{}).Where("incident_id = ?", incidentID)
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.Investigation]{}, domain.ErrInvalidInvestigation.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.Investigation]{}, mapError(err)
	}
	return paging.Slice[domain.Investigation]{Items: items, Total: total}, nil
}
