package mapper

import (
	githubdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/github"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func GitHubRepositoryToDTO(value domain.GitHubRepository) githubdto.GitHubRepositoryDTO {
	return githubdto.GitHubRepositoryDTO{
		GitHubAccountID: value.GitHubAccountID.String(), RepositoryIdentifier: value.RepositoryIdentifier,
		Owner: value.Owner, Name: value.Name, FullName: value.FullName, DefaultBranch: value.DefaultBranch,
		Private: value.Private, HTMLURL: value.HTMLURL,
	}
}
