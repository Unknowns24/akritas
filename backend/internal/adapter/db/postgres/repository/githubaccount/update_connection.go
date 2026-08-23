package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) UpdateConnection(ctx context.Context, account *domain.GitHubAccount) error {
	result := r.db.WithContext(ctx).Table("github_accounts").Model(&domain.GitHubAccount{}).Where("id = ?", account.ID).Updates(map[string]any{
		"authentication_status": account.AuthenticationStatus,
		"repository_count":      account.RepositoryCount,
		"last_checked_at":       account.LastCheckedAt,
		"updated_at":            account.UpdatedAt,
	})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrIntegrationNotFound
	}
	return nil
}
