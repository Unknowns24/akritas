package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func UpdateGitHubAccountToCommand(value dto.UpdateGitHubAccountRequestDTO) portsin.UpdateGitHubAccountCommand {
	return portsin.UpdateGitHubAccountCommand{DisplayName: value.DisplayName, PersonalAccessToken: value.PersonalAccessToken}
}
