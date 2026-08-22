package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) CountByGitHubAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&domain.Project{}).Where("github_account_id = ?", accountID).Count(&total).Error
	return total, err
}
