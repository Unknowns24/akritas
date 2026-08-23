package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) GetMonitoring(ctx context.Context, id uuid.UUID) (domain.MonitoringConfiguration, error) {
	project, err := uc.projects.Get(ctx, id)
	if err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	return project.MonitoringConfiguration.Clone(), nil
}
