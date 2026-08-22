package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	inUse, err := uc.usage.GitHubAccountInUse(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return domain.ErrIntegrationInUse
	}
	return uc.store.Delete(ctx, id)
}
