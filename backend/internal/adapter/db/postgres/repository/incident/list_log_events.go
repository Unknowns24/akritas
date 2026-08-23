package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) ListLogEvents(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.LogEvent], error) {
	base := txcontext.From(ctx, r.db).WithContext(ctx).Table("log_events").Model(&domain.LogEvent{}).Where("incident_id = ?", incidentID)
	query, err := ukerpagination.Apply(base, params)
	if err != nil {
		return paging.Slice[domain.LogEvent]{}, domain.ErrInvalidLogEvent.Wrap(err)
	}
	var items []domain.LogEvent
	if err := query.Find(&items).Error; err != nil {
		return paging.Slice[domain.LogEvent]{}, mapError(err)
	}
	countBase := txcontext.From(ctx, r.db).WithContext(ctx).Table("log_events").Model(&domain.LogEvent{}).Where("incident_id = ?", incidentID)
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.LogEvent]{}, domain.ErrInvalidLogEvent.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.LogEvent]{}, mapError(err)
	}
	return paging.Slice[domain.LogEvent]{Items: items, Total: total}, nil
}
