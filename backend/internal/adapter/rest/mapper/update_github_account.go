package mapper

import (
	githubdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/github"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func UpdateGitHubAccountToCommand(value githubdto.UpdateGitHubAccountRequestDTO) portsin.UpdateGitHubAccountCommand {
	return portsin.UpdateGitHubAccountCommand{DisplayName: value.DisplayName, PersonalAccessToken: value.PersonalAccessToken}
}
