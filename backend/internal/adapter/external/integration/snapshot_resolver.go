package integration

import (
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type SnapshotResolver struct{}

func NewSnapshotResolver() SnapshotResolver {
	return SnapshotResolver{}
}

func (SnapshotResolver) ResolveGitHubRepository(account *domain.GitHubAccount, repositoryIdentifier, defaultBranch string) (domain.GitHubRepository, error) {
	if account == nil {
		return domain.GitHubRepository{}, apperr.ErrGitHubAccountNotFound
	}
	identifier := strings.TrimSpace(repositoryIdentifier)
	if identifier == "" || strings.TrimSpace(defaultBranch) == "" {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable
	}
	owner, name := strings.TrimSpace(account.AccountIdentifier), identifier
	if parts := strings.Split(identifier, "/"); len(parts) == 2 {
		owner, name = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	} else if strings.Contains(identifier, "/") {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable
	}
	if owner == "" || name == "" {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable
	}
	repository, err := domain.NewGitHubRepository(
		account.ID, identifier, owner, name, defaultBranch, false,
		"https://github.com/"+owner+"/"+name,
	)
	if err != nil {
		return domain.GitHubRepository{}, apperr.ErrRepositoryNotResolvable.Wrap(err)
	}
	return repository, nil
}

func (SnapshotResolver) ResolveDokployApplication(server *domain.DokployServer, applicationIdentifier string) (domain.DokployApplication, error) {
	if server == nil {
		return domain.DokployApplication{}, apperr.ErrDokployServerNotFound
	}
	identifier := strings.TrimSpace(applicationIdentifier)
	if identifier == "" {
		return domain.DokployApplication{}, apperr.ErrApplicationNotResolvable
	}
	application, err := domain.NewDokployApplication(
		server.ID, identifier, identifier, identifier, "", domain.DokployApplicationUnknown,
	)
	if err != nil {
		return domain.DokployApplication{}, apperr.ErrApplicationNotResolvable.Wrap(err)
	}
	return application, nil
}
