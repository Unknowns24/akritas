package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"gorm.io/gorm"
)

func (r *Repository) Update(ctx context.Context, account *domain.GitHubAccount, secret *portsout.SecretValue) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table("github_accounts").Model(&domain.GitHubAccount{}).Where("id = ?", account.ID).Select("display_name", "account_type", "authentication_method", "account_identifier", "authentication_status", "credential_configured", "repository_count", "last_checked_at", "manage_url", "updated_at").Updates(account)
		if result.Error != nil {
			return mapError(result.Error)
		}
		if result.RowsAffected == 0 {
			return domain.ErrIntegrationNotFound
		}
		if secret != nil {
			return r.credentials.PutTx(ctx, tx, portsout.CredentialOwnerGitHubAccount, account.ID, *secret)
		}
		return nil
	})
}
