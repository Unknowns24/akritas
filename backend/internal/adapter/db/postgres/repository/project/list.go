package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) List(ctx context.Context, params paging.Params) (paging.Slice[domain.Project], error) {
	db := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Model(&domain.Project{})
	query, err := ukerpagination.Apply(db, params)
	if err != nil {
		return paging.Slice[domain.Project]{}, domain.ErrInvalidProject.Wrap(err)
	}
	var items []domain.Project
	if err := query.Find(&items).Error; err != nil {
		return paging.Slice[domain.Project]{}, mapError(err)
	}
	for i := range items {
		if err := items[i].Validate(); err != nil {
			return paging.Slice[domain.Project]{}, mapError(err)
		}
	}
	countBase := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Model(&domain.Project{})
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.Project]{}, domain.ErrInvalidProject.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.Project]{}, mapError(err)
	}
	return paging.Slice[domain.Project]{Items: items, Total: total}, nil
}
