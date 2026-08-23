package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) List(ctx context.Context, params paging.Params) (paging.Slice[domain.DokployServer], error) {
	dataBase := r.db.WithContext(ctx).Table("dokploy_servers").Model(&domain.DokployServer{})
	dataQuery, err := ukerpagination.Apply(dataBase, params)
	if err != nil {
		return paging.Slice[domain.DokployServer]{}, domain.ErrInvalidDokployServer.Wrap(err)
	}
	var items []domain.DokployServer
	if err := dataQuery.Find(&items).Error; err != nil {
		return paging.Slice[domain.DokployServer]{}, mapError(err)
	}
	for i := range items {
		if err := items[i].Validate(); err != nil {
			return paging.Slice[domain.DokployServer]{}, mapError(err)
		}
	}
	countBase := r.db.WithContext(ctx).Table("dokploy_servers").Model(&domain.DokployServer{})
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.DokployServer]{}, domain.ErrInvalidDokployServer.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.DokployServer]{}, mapError(err)
	}
	return paging.Slice[domain.DokployServer]{Items: items, Total: total}, nil
}
