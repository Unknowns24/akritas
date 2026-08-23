package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var value domain.Project
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("id = ?", id).First(&value).Error; err != nil {
		return nil, mapError(err)
	}
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
