package githubapp

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"gorm.io/gorm"
)

func (r *Repository) CompleteConversion(ctx context.Context, registration *portsout.GitHubAppRegistration, secrets []portsout.SecretValue) error {
	if registration == nil || registration.Status != portsout.GitHubAppRegistrationConverted || len(registration.InstallationStateDigest) != 32 || registration.AppID == nil || *registration.AppID < 1 || len(secrets) != 2 {
		return domain.ErrManifestStateInvalid
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table("github_app_registrations").Model(&domain.GitHubAppRegistration{}).Where("id = ? AND status = ? AND conversion_consumed_at IS NOT NULL", registration.ID, portsout.GitHubAppRegistrationCreated).Updates(registration)
		if result.Error != nil {
			return mapError(result.Error)
		}
		if result.RowsAffected != 1 {
			return domain.ErrManifestStateConflict
		}
		for _, secret := range secrets {
			if err := r.credentials.PutTx(ctx, tx, portsout.CredentialOwnerGitHubManifest, registration.ID, secret); err != nil {
				return err
			}
		}
		return nil
	})
}
