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
	if err := uc.ensureApplicationAvailable(ctx, command.DokployServerID, command.ApplicationIdentifier, uuid.Nil); err != nil {
		return nil, err
	}
	repository, err := uc.resolveRepository(ctx, command.GitHubAccountID, command.RepositoryIdentifier, command.DefaultBranch)
	if err != nil {
		return nil, err
	}
	application, err := uc.resolveApplication(ctx, command.DokployServerID, command.ApplicationIdentifier)
	if err != nil {
		return nil, err
	}
	created, err := domain.NewProject(uc.newID(), command.Name, command.Description, repository, application, command.MonitoringConfiguration, uc.now().UTC())
	if err != nil {
		return nil, err
	}
	if err := uc.projects.Create(ctx, created); err != nil {
		return nil, err
	}
	return result(created), nil
}
