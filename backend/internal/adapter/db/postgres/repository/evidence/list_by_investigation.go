package evidence

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) ListByInvestigation(ctx context.Context, investigationID uuid.UUID, params paging.Params) (paging.Slice[domain.Evidence], error) {
	db := txcontext.From(ctx, r.db).WithContext(ctx).Table("evidence").Model(&domain.Evidence{}).Where("investigation_id = ?", investigationID)
	query, err := ukerpagination.Apply(db, params)
	if err != nil {
		return paging.Slice[domain.Evidence]{}, domain.ErrInvalidEvidence.Wrap(err)
	}
	var items []domain.Evidence
	if err := query.Find(&items).Error; err != nil {
		return paging.Slice[domain.Evidence]{}, mapError(err)
	}
	for i := range items {
		if err := items[i].Validate(); err != nil {
			return paging.Slice[domain.Evidence]{}, mapError(err)
		}
	}
	countBase := txcontext.From(ctx, r.db).WithContext(ctx).Table("evidence").Model(&domain.Evidence{}).Where("investigation_id = ?", investigationID)
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.Evidence]{}, domain.ErrInvalidEvidence.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.Evidence]{}, mapError(err)
	}
	return paging.Slice[domain.Evidence]{Items: items, Total: total}, nil
}
