package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/google/uuid"
)

func (uc *UseCase) Update(ctx context.Context, command inproject.UpdateCommand) (*inproject.Result, error) {
	if command.ID == uuid.Nil {
		return nil, apperr.ErrInvalidProjectCommand
	}
	project, err := uc.projects.GetByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}
	now := uc.now()
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
	if command.GitHubAccountID != nil || command.RepositoryIdentifier != nil || command.DefaultBranch != nil ||
		command.DokployServerID != nil || command.ApplicationIdentifier != nil {
		accountID := project.GitHubRepository.GitHubAccountID
		repositoryIdentifier := project.GitHubRepository.RepositoryIdentifier
		defaultBranch := project.GitHubRepository.DefaultBranch
		serverID := project.DokployApplication.DokployServerID
		applicationIdentifier := project.DokployApplication.ApplicationIdentifier
		if command.GitHubAccountID != nil {
			accountID = *command.GitHubAccountID
		}
		if command.RepositoryIdentifier != nil {
			repositoryIdentifier = *command.RepositoryIdentifier
		}
		if command.DefaultBranch != nil {
			defaultBranch = *command.DefaultBranch
		}
		if command.DokployServerID != nil {
			serverID = *command.DokployServerID
		}
		if command.ApplicationIdentifier != nil {
			applicationIdentifier = *command.ApplicationIdentifier
		}
		account, err := uc.accounts.GetByID(ctx, accountID)
		if err != nil {
			return nil, err
		}
		server, err := uc.servers.GetByID(ctx, serverID)
		if err != nil {
			return nil, err
		}
		if err := uc.ensureApplicationAvailable(ctx, server.ID, applicationIdentifier, project.ID); err != nil {
			return nil, err
		}
		repository, err := uc.snapshots.ResolveGitHubRepository(account, repositoryIdentifier, defaultBranch)
		if err != nil {
			return nil, err
		}
		application, err := uc.snapshots.ResolveDokployApplication(server, applicationIdentifier)
		if err != nil {
			return nil, err
		}
		if err := project.ReplaceIntegrations(repository, application, now); err != nil {
			return nil, err
		}
	}
	if err := uc.projects.Update(ctx, project); err != nil {
		return nil, err
	}
	return inproject.NewResult(project), nil
}
