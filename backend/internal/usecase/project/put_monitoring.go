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
	if uc.monitoring == nil || uc.transactor == nil {
		if err := uc.projects.Update(ctx, project, expected); err != nil {
			return domain.MonitoringConfiguration{}, err
		}
	} else if err := uc.transactor.WithinTransaction(ctx, func(txctx context.Context) error {
		if updateErr := uc.projects.Update(txctx, project, expected); updateErr != nil {
			return updateErr
		}
		if !configuration.Enabled {
			return nil
		}
		checkpoint, getErr := uc.monitoring.GetCurrentCheckpoint(txctx, id, true)
		if getErr != nil {
			return getErr
		}
		if checkpoint == nil {
			created, createErr := domain.NewMonitoringCheckpoint(uc.newID(), *project, domain.InitialLogIngestionFromNow, now)
			if createErr != nil {
				return createErr
			}
			return uc.monitoring.CreateCheckpoint(txctx, created)
		}
		if checkpoint.CursorTimestamp != nil || checkpoint.InitialBackfillPending {
			return nil
		}
		version := checkpoint.Version
		anchor := now
		checkpoint.CursorTimestamp = &anchor
		checkpoint.CursorContentHash = "anchor"
		checkpoint.Version++
		checkpoint.UpdatedAt = now
		return uc.monitoring.UpdateCheckpoint(txctx, checkpoint, version)
	}); err != nil {
		return domain.MonitoringConfiguration{}, err
	}
	return project.MonitoringConfiguration.Clone(), nil
}
