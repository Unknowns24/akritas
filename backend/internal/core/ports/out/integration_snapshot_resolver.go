package out

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

type IntegrationSnapshotResolver interface {
	ResolveGitHubRepository(account *domain.GitHubAccount, repositoryIdentifier, defaultBranch string) (domain.GitHubRepository, error)
	ResolveDokployApplication(server *domain.DokployServer, applicationIdentifier string) (domain.DokployApplication, error)
}
