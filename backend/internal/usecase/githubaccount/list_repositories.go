package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (uc *UseCase) ListRepositories(ctx context.Context, id uuid.UUID, params paging.Params) (paging.Slice[domain.GitHubRepository], error) {
	account, err := uc.store.Get(ctx, id)
	if err != nil {
		return paging.Slice[domain.GitHubRepository]{}, err
	}
	page, err := uc.gateway.ListRepositories(ctx, *account, params)
	if err != nil {
		return paging.Slice[domain.GitHubRepository]{}, err
	}
	now := uc.now().UTC()
	account.RepositoryCount = int(page.Total)
	account.LastCheckedAt = &now
	account.UpdatedAt = now
	if err := uc.store.UpdateConnection(ctx, account); err != nil {
		return paging.Slice[domain.GitHubRepository]{}, err
	}
	return page, nil
}
