package project

import (
	"context"
	"errors"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"gorm.io/gorm"
)

func (r *Repository) GetByNormalizedName(ctx context.Context, name string) (*domain.Project, error) {
	var project domain.Project
	err := r.db.WithContext(ctx).Where("LOWER(name) = ?", strings.ToLower(strings.TrimSpace(name))).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrProjectNotFound
		}
		return nil, err
	}
	if err := project.Validate(); err != nil {
		return nil, err
	}
	return &project, nil
}
