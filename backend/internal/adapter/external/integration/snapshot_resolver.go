package integration

import (
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// SnapshotResolver fills non-secret GitHub/Dokploy projections from identifiers.
// It does not call GitHub or Dokploy APIs; a later adapter behind
// ports/out.IntegrationSnapshotResolver can replace this with live metadata
// once Credential Store exists.
//
// MVP host is github.com only (ADR-009). private is always false because
// visibility is unknown without the GitHub API. html_url is synthesized as
// https://github.com/{owner}/{name}. No PAT, token, or private key is stored
// on the GitHubRepository value object.
//
// Dokploy: instance_identifier and display_name default to application_identifier.
// environment is empty and status is unknown because those are unknown without
// the Dokploy API. No API key or credential is stored on the DokployApplication
// value object. An empty identifier is 0x503005N (ErrApplicationNotResolvable).
type SnapshotResolver struct{}

func NewSnapshotResolver() SnapshotResolver {
	return SnapshotResolver{}
}

func (SnapshotResolver) ResolveGitHubRepository(account *domain.GitHubAccount, repositoryIdentifier, defaultBranch string) (domain.GitHubRepository, error) {
	if account == nil {
		return domain.GitHubRepository{}, apperr.ErrGitHubAccountNotFound
	}
	owner, name, identifier, err := parseGitHubIdentifier(account.AccountIdentifier, repositoryIdentifier)
	if err != nil || strings.TrimSpace(defaultBranch) == "" {
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

func parseGitHubIdentifier(accountIdentifier, repositoryIdentifier string) (owner, name, identifier string, err error) {
	identifier = strings.TrimSpace(repositoryIdentifier)
	if identifier == "" {
		return "", "", "", apperr.ErrRepositoryNotResolvable
	}
	owner, name = strings.TrimSpace(accountIdentifier), identifier
	if parts := strings.Split(identifier, "/"); len(parts) == 2 {
		owner, name = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	} else if strings.Contains(identifier, "/") {
		return "", "", "", apperr.ErrRepositoryNotResolvable
	}
	if owner == "" || name == "" {
		return "", "", "", apperr.ErrRepositoryNotResolvable
	}
	return owner, name, identifier, nil
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
