package project

import (
	"context"
	"errors"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetByDokployApplication(ctx context.Context, serverID uuid.UUID, applicationIdentifier string) (*domain.Project, error) {
	var project domain.Project
	err := r.db.WithContext(ctx).
		Where("dokploy_server_id = ? AND application_identifier = ?", serverID, strings.TrimSpace(applicationIdentifier)).
		First(&project).Error
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
