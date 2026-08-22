package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (r *Repository) List(ctx context.Context, params paging.Params) (paging.Slice[domain.GitHubAccount], error) {
	dataBase := r.db.WithContext(ctx).Table("github_accounts").Model(&domain.GitHubAccount{})
	dataQuery, err := ukerpagination.Apply(dataBase, params)
	if err != nil {
		return paging.Slice[domain.GitHubAccount]{}, domain.ErrInvalidGitHubAccount.Wrap(err)
	}
	var items []domain.GitHubAccount
	if err := dataQuery.Find(&items).Error; err != nil {
		return paging.Slice[domain.GitHubAccount]{}, mapError(err)
	}
	for i := range items {
		if err := items[i].Validate(); err != nil {
			return paging.Slice[domain.GitHubAccount]{}, mapError(err)
		}
	}
	countBase := r.db.WithContext(ctx).Table("github_accounts").Model(&domain.GitHubAccount{})
	countQuery, err := ukerpagination.ApplyFilters(countBase, params.Filters)
	if err != nil {
		return paging.Slice[domain.GitHubAccount]{}, domain.ErrInvalidGitHubAccount.Wrap(err)
	}
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return paging.Slice[domain.GitHubAccount]{}, mapError(err)
	}
	return paging.Slice[domain.GitHubAccount]{Items: items, Total: total}, nil
}
