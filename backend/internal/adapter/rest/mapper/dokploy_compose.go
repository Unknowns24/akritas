package mapper

import (
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokployComposeToDTO(value domain.DokployCompose) dokploydto.DokployComposeDTO {
	return dokploydto.DokployComposeDTO{DokployServerID: value.DokployServerID.String(), ComposeIdentifier: value.ComposeIdentifier,
		InstanceIdentifier: value.InstanceIdentifier, DisplayName: value.DisplayName,
		EnvironmentIdentifier: value.EnvironmentIdentifier, Status: string(value.Status)}
}
