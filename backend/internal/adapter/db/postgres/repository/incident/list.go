package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) List(ctx context.Context, params paging.Params) (paging.Slice[domain.Incident], error) {
	base := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").
		Select("incidents.*, (SELECT name FROM projects WHERE projects.id = incidents.project_id) AS project_name")
	query, err := ukerpagination.Apply(base, params)
	if err != nil {
		return paging.Slice[domain.Incident]{}, domain.ErrInvalidIncident.Wrap(err)
	}
	var rows []incidentRow
	if err := query.Scan(&rows).Error; err != nil {
		return paging.Slice[domain.Incident]{}, mapError(err)
	}
	items := make([]domain.Incident, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.domain())
	}
	countBase := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents")
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.Incident]{}, domain.ErrInvalidIncident.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.Incident]{}, mapError(err)
	}
	return paging.Slice[domain.Incident]{Items: items, Total: total}, nil
}
