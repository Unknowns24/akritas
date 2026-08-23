package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*domain.GitHubAccount, error) {
	return uc.store.Get(ctx, id)
}
