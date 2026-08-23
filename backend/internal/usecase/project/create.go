package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func (uc *UseCase) Create(ctx context.Context, command portsin.CreateProjectCommand) (*portsin.ProjectResult, error) {
	if err := uc.ensureNameAvailable(ctx, command.Name, uuid.Nil); err != nil {
		return nil, err
	}
	if err := uc.ensureSourceAvailable(ctx, command.DokploySource, uuid.Nil); err != nil {
		return nil, err
	}
	repository, err := uc.resolveRepository(ctx, command.GitHubAccountID, command.RepositoryIdentifier, command.DefaultBranch)
	if err != nil {
		return nil, err
	}
	source, err := uc.resolveSource(ctx, command.DokploySource)
	if err != nil {
		return nil, err
	}
	created, err := domain.NewProject(uc.newID(), command.Name, command.Description, repository, source, command.MonitoringConfiguration, uc.now().UTC())
	if err != nil {
		return nil, err
	}
	ingestion := command.InitialLogIngestion
	if ingestion == "" {
		ingestion = domain.InitialLogIngestionFromNow
	}
	if err := ingestion.Validate(); err != nil {
		return nil, err
	}
	if uc.monitoring == nil || uc.transactor == nil {
		if err := uc.projects.Create(ctx, created); err != nil {
			return nil, err
		}
	} else {
		checkpoint, checkpointErr := domain.NewMonitoringCheckpoint(uc.newID(), *created, ingestion, created.CreatedAt)
		if checkpointErr != nil {
			return nil, checkpointErr
		}
		if err := uc.transactor.WithinTransaction(ctx, func(txctx context.Context) error {
			if createErr := uc.projects.Create(txctx, created); createErr != nil {
				return createErr
			}
			return uc.monitoring.CreateCheckpoint(txctx, checkpoint)
		}); err != nil {
			return nil, err
		}
	}
	return result(created), nil
}
