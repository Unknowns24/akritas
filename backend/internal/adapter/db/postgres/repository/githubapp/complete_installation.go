package githubapp

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"gorm.io/gorm"
)

func (r *Repository) CompleteInstallation(ctx context.Context, registration *portsout.GitHubAppRegistration, account *domain.GitHubAccount, binding portsout.GitHubAppBinding) error {
	if registration == nil || account == nil || binding.GitHubAccountID != account.ID || registration.Status != portsout.GitHubAppRegistrationCompleted {
		return domain.ErrManifestStateInvalid
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("github_accounts").Create(account).Error; err != nil {
			return mapError(err)
		}
		binding.CreatedAt = account.CreatedAt
		binding.UpdatedAt = account.UpdatedAt
		if err := tx.Table("github_app_bindings").Create(&binding).Error; err != nil {
			return mapError(err)
		}
		if err := r.credentials.MoveOwnerTx(ctx, tx, portsout.CredentialOwnerGitHubManifest, registration.ID, portsout.CredentialOwnerGitHubAccount, account.ID); err != nil {
			return err
		}
		result := tx.Table("github_app_registrations").Model(&domain.GitHubAppRegistration{}).Where("id = ? AND status = ? AND installation_consumed_at IS NOT NULL", registration.ID, portsout.GitHubAppRegistrationConverted).Updates(registration)
		if result.Error != nil {
			return mapError(result.Error)
		}
		if result.RowsAffected != 1 {
			return domain.ErrManifestStateConflict
		}
		return nil
	})
}
