package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func GitHubManifestRegistrationToCommand(value dto.GitHubManifestRegistrationRequestDTO) portsin.StartGitHubAppRegistrationCommand {
	return portsin.StartGitHubAppRegistrationCommand{DisplayName: value.DisplayName, OwnerType: domain.GitHubAccountType(value.OwnerType), Organization: value.Organization}
}
