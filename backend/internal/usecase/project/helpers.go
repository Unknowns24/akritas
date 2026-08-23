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

func (uc *UseCase) resolveApplication(ctx context.Context, serverID uuid.UUID, identifier string) (domain.DokployApplication, error) {
	server, err := uc.servers.Get(ctx, serverID)
	if err != nil {
		return domain.DokployApplication{}, err
	}
	application, err := uc.dokploy.GetApplication(ctx, *server, strings.TrimSpace(identifier))
	if errors.Is(err, domain.ErrIntegrationNotFound) {
		return domain.DokployApplication{}, domain.ErrProjectApplicationNotFound
	}
	return application, err
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

func (uc *UseCase) ensureApplicationAvailable(ctx context.Context, serverID uuid.UUID, identifier string, exclude uuid.UUID) error {
	existing, err := uc.projects.FindByDokployApplication(ctx, serverID, strings.TrimSpace(identifier))
	if errors.Is(err, domain.ErrProjectNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if existing.ID != exclude {
		return domain.ErrProjectApplicationConflict
	}
	return nil
}
