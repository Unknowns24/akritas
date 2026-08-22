package githubaccount

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GitHubAccount, error) {
	var account domain.GitHubAccount
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrGitHubAccountNotFound
		}
		return nil, err
	}
	if err := account.Validate(); err != nil {
		return nil, err
	}
	return &account, nil
}
