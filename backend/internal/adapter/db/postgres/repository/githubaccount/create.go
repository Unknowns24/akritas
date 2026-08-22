package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Create(ctx context.Context, account *domain.GitHubAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}
