package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) PutMonitoring(ctx context.Context, id uuid.UUID, configuration domain.MonitoringConfiguration) (domain.MonitoringConfiguration, error) {
	project, err := uc.projects.Get(ctx, id)
	if err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	expected := project.UpdatedAt
	now := uc.now().UTC()
	if configuration.Enabled {
		repository, resolveErr := uc.resolveRepository(ctx, project.GitHubRepository.GitHubAccountID, project.GitHubRepository.RepositoryIdentifier, project.GitHubRepository.DefaultBranch)
		if resolveErr != nil {
			return domain.MonitoringConfiguration{}, resolveErr
		}
		application, resolveErr := uc.resolveApplication(ctx, project.DokployApplication.DokployServerID, project.DokployApplication.ApplicationIdentifier)
		if resolveErr != nil {
			return domain.MonitoringConfiguration{}, resolveErr
		}
		if refreshErr := project.RefreshIntegrationSnapshots(repository, application, now); refreshErr != nil {
			return domain.MonitoringConfiguration{}, refreshErr
		}
	}
	if err := project.ReplaceMonitoringConfiguration(configuration, now); err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	if err := uc.projects.Update(ctx, project, expected); err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	return project.MonitoringConfiguration.Clone(), nil
}
