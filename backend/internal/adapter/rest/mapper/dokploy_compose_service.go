package mapper

import (
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokployComposeServiceToDTO(value domain.DokployComposeService) dokploydto.DokployComposeServiceDTO {
	return dokploydto.DokployComposeServiceDTO{DokployServerID: value.DokployServerID.String(), ComposeIdentifier: value.ComposeIdentifier, ServiceName: value.ServiceName}
}
