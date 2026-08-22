package project

import (
	"context"
	"errors"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/google/uuid"
)

func (uc *UseCase) Create(ctx context.Context, command inproject.CreateCommand) (*inproject.Result, error) {
	account, err := uc.accounts.GetByID(ctx, command.GitHubAccountID)
	if err != nil {
		return nil, err
	}
	server, err := uc.servers.GetByID(ctx, command.DokployServerID)
	if err != nil {
		return nil, err
	}
	repository, err := uc.snapshots.ResolveGitHubRepository(account, command.RepositoryIdentifier, command.DefaultBranch)
	if err != nil {
		return nil, err
	}
	application, err := uc.snapshots.ResolveDokployApplication(server, command.ApplicationIdentifier)
	if err != nil {
		return nil, err
	}
	if err := uc.ensureNameAvailable(ctx, command.Name, uuid.Nil); err != nil {
		return nil, err
	}
	if err := uc.ensureApplicationAvailable(ctx, server.ID, command.ApplicationIdentifier, uuid.Nil); err != nil {
		return nil, err
	}
	now := uc.now()
	created, err := domain.NewProject(uuid.New(), command.Name, command.Description, repository, application, command.MonitoringConfiguration, now)
	if err != nil {
		return nil, err
	}
	if err := uc.projects.Create(ctx, created); err != nil {
		return nil, err
	}
	return inproject.NewResult(created), nil
}

func (uc *UseCase) ensureNameAvailable(ctx context.Context, name string, exclude uuid.UUID) error {
	existing, err := uc.projects.GetByNormalizedName(ctx, strings.TrimSpace(name))
	if err != nil {
		if errors.Is(err, apperr.ErrProjectNotFound) {
			return nil
		}
		return err
	}
	if existing.ID != exclude {
		return apperr.ErrProjectNameConflict
	}
	return nil
}

func (uc *UseCase) ensureApplicationAvailable(ctx context.Context, serverID uuid.UUID, applicationIdentifier string, exclude uuid.UUID) error {
	existing, err := uc.projects.GetByDokployApplication(ctx, serverID, strings.TrimSpace(applicationIdentifier))
	if err != nil {
		if errors.Is(err, apperr.ErrProjectNotFound) {
			return nil
		}
		return err
	}
	if existing.ID != exclude {
		return apperr.ErrProjectApplicationConflict
	}
	return nil
}
