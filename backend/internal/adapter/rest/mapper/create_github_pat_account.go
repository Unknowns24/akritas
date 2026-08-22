package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func CreateGitHubPATAccountToCommand(value dto.CreateGitHubPATAccountRequestDTO) portsin.CreateGitHubPATAccountCommand {
	return portsin.CreateGitHubPATAccountCommand{DisplayName: value.DisplayName, AccountType: domain.GitHubAccountType(value.AccountType), AccountIdentifier: value.AccountIdentifier, PersonalAccessToken: value.PersonalAccessToken}
}
