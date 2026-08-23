package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func (uc *UseCase) Update(ctx context.Context, command portsin.UpdateProjectCommand) (*portsin.ProjectResult, error) {
	project, err := uc.projects.Get(ctx, command.ID)
	if err != nil {
		return nil, err
	}
	expected := project.UpdatedAt
	now := uc.now().UTC()
	if command.Name != nil || command.Description != nil {
		name, description := project.Name, project.Description
		if command.Name != nil {
			name = *command.Name
		}
		if command.Description != nil {
			description = *command.Description
		}
		if err := uc.ensureNameAvailable(ctx, name, project.ID); err != nil {
			return nil, err
		}
		if err := project.Rename(name, description, now); err != nil {
			return nil, err
		}
	}
	githubChanged := command.GitHubAccountID != nil || command.RepositoryIdentifier != nil || command.DefaultBranch != nil
	dokployChanged := command.DokployServerID != nil || command.ApplicationIdentifier != nil
	if githubChanged || dokployChanged {
		if project.MonitoringConfiguration.Enabled {
			return nil, domain.ErrProjectMustBeDisabled
		}
		repository, application := project.GitHubRepository, project.DokployApplication
		if githubChanged {
			accountID, identifier, branch := repository.GitHubAccountID, repository.RepositoryIdentifier, repository.DefaultBranch
			if command.GitHubAccountID != nil {
				accountID = *command.GitHubAccountID
			}
			if command.RepositoryIdentifier != nil {
				identifier = *command.RepositoryIdentifier
			}
			if command.DefaultBranch != nil {
				branch = *command.DefaultBranch
			}
			repository, err = uc.resolveRepository(ctx, accountID, identifier, branch)
			if err != nil {
				return nil, err
			}
		}
		if dokployChanged {
			serverID, identifier := application.DokployServerID, application.ApplicationIdentifier
			if command.DokployServerID != nil {
				serverID = *command.DokployServerID
			}
			if command.ApplicationIdentifier != nil {
				identifier = *command.ApplicationIdentifier
			}
			if err := uc.ensureApplicationAvailable(ctx, serverID, identifier, project.ID); err != nil {
				return nil, err
			}
			application, err = uc.resolveApplication(ctx, serverID, identifier)
			if err != nil {
				return nil, err
			}
		}
		if err := project.ReplaceIntegrations(repository, application, now); err != nil {
			return nil, err
		}
	}
	if err := uc.projects.Update(ctx, project, expected); err != nil {
		return nil, err
	}
	return result(project), nil
}
