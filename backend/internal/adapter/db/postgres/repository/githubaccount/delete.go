package githubaccount

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.credentials.DeleteOwnerTx(ctx, tx, portsout.CredentialOwnerGitHubAccount, id); err != nil {
			return err
		}
		result := tx.Table("github_accounts").Delete(&domain.GitHubAccount{}, "id = ?", id)
		if result.Error != nil {
			return mapError(result.Error)
		}
		if result.RowsAffected == 0 {
			return domain.ErrIntegrationNotFound
		}
		return nil
	})
}
