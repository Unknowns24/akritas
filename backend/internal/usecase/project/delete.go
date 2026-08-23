package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	project, err := uc.projects.Get(ctx, id)
	if err != nil {
		return err
	}
	if project.MonitoringConfiguration.Enabled {
		return domain.ErrProjectMustBeDisabled
	}
	return uc.projects.Delete(ctx, id)
}
