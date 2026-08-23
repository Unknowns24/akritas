package githubaccount

import (
	"context"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*domain.GitHubAccount, error) {
	var account domain.GitHubAccount
	if err := r.db.WithContext(ctx).Table("github_accounts").First(&account, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	if err := account.Validate(); err != nil {
		return nil, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return &account, nil
}
