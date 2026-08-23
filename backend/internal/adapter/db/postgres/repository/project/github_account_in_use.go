package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/google/uuid"
)

func (r *Repository) GitHubAccountInUse(ctx context.Context, id uuid.UUID) (bool, error) {
	var total int64
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("github_account_id = ?", id).Count(&total).Error
	if err != nil {
		return false, mapError(err)
	}
	return total > 0, nil
}
