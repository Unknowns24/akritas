package project

import (
	"context"
	"errors"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func result(value *domain.Project) *portsin.ProjectResult {
	return &portsin.ProjectResult{Project: value, BuiltInDetectionRules: domain.AllBuiltInDetectionRules()}
}

func (uc *UseCase) resolveRepository(ctx context.Context, accountID uuid.UUID, identifier, expectedBranch string) (domain.GitHubRepository, error) {
	account, err := uc.accounts.Get(ctx, accountID)
	if err != nil {
		return domain.GitHubRepository{}, err
	}
	repository, err := uc.github.GetRepository(ctx, *account, strings.TrimSpace(identifier))
	if errors.Is(err, domain.ErrIntegrationNotFound) {
		return domain.GitHubRepository{}, domain.ErrProjectRepositoryNotFound
	}
	if err != nil {
		return domain.GitHubRepository{}, err
	}
	if repository.DefaultBranch != strings.TrimSpace(expectedBranch) {
		return domain.GitHubRepository{}, domain.ErrProjectDefaultBranchMismatch
	}
	return repository, nil
}

func (uc *UseCase) resolveSource(ctx context.Context, selector domain.DokploySourceSelector) (domain.DokploySource, error) {
	selector = selector.Normalize()
	if err := selector.Validate(); err != nil {
		return domain.DokploySource{}, err
	}
	server, err := uc.servers.Get(ctx, selector.DokployServerID)
	if err != nil {
		return domain.DokploySource{}, err
	}
	if selector.Type == domain.DokploySourceApplication {
		application, resolveErr := uc.dokploy.GetApplication(ctx, *server, selector.ResourceIdentifier)
		if errors.Is(resolveErr, domain.ErrIntegrationNotFound) {
			return domain.DokploySource{}, domain.ErrProjectDokploySourceNotFound
		}
		if resolveErr != nil {
			return domain.DokploySource{}, resolveErr
		}
		return domain.SourceFromApplication(application)
	}
	compose, err := uc.dokploy.GetCompose(ctx, *server, selector.ResourceIdentifier)
	if errors.Is(err, domain.ErrIntegrationNotFound) {
		return domain.DokploySource{}, domain.ErrProjectDokploySourceNotFound
	}
	if err != nil {
		return domain.DokploySource{}, err
	}
	services, err := uc.dokploy.ListComposeServices(ctx, *server, selector.ResourceIdentifier, false)
	if err != nil {
		return domain.DokploySource{}, err
	}
	for _, service := range services {
		if service.ServiceName == selector.ServiceName {
			return compose.Source(selector.ServiceName)
		}
	}
	return domain.DokploySource{}, domain.ErrProjectDokploySourceNotFound
}

func (uc *UseCase) ensureNameAvailable(ctx context.Context, name string, exclude uuid.UUID) error {
	existing, err := uc.projects.FindByNormalizedName(ctx, strings.TrimSpace(name))
	if errors.Is(err, domain.ErrProjectNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if existing.ID != exclude {
		return domain.ErrProjectNameConflict
	}
	return nil
}

func (uc *UseCase) ensureSourceAvailable(ctx context.Context, selector domain.DokploySourceSelector, exclude uuid.UUID) error {
	selector = selector.Normalize()
	if err := selector.Validate(); err != nil {
		return err
	}
	existing, err := uc.projects.FindByDokploySource(ctx, selector)
	if errors.Is(err, domain.ErrProjectNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if existing.ID != exclude {
		return domain.ErrProjectDokploySourceConflict
	}
	return nil
}
