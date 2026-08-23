package githubapp

import (
	"context"
	"errors"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetBinding(ctx context.Context, accountID uuid.UUID) (portsout.GitHubAppBinding, error) {
	var binding domain.GitHubAppBinding
	if err := r.db.WithContext(ctx).Table("github_app_bindings").First(&binding, "github_account_id = ?", accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return portsout.GitHubAppBinding{}, domain.ErrIntegrationNotFound
		}
		return portsout.GitHubAppBinding{}, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return binding, nil
}
