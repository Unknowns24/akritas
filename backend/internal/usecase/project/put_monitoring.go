package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
)

func (uc *UseCase) PutMonitoring(ctx context.Context, command inproject.MonitoringCommand) (domain.MonitoringConfiguration, error) {
	project, err := uc.projects.GetByID(ctx, command.ProjectID)
	if err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	if err := project.ReplaceMonitoringConfiguration(command.MonitoringConfiguration, uc.now()); err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	if err := uc.projects.Update(ctx, project); err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	return project.MonitoringConfiguration, nil
}
