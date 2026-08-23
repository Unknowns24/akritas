package project

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) FindByNormalizedName(ctx context.Context, name string) (*domain.Project, error) {
	var value domain.Project
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("LOWER(name) = ?", strings.ToLower(strings.TrimSpace(name))).First(&value).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
