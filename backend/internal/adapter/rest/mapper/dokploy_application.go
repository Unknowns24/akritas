package mapper

import (
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokployApplicationToDTO(value domain.DokployApplication) dokploydto.DokployApplicationDTO {
	return dokploydto.DokployApplicationDTO{
		DokployServerID: value.DokployServerID.String(), ApplicationIdentifier: value.ApplicationIdentifier,
		InstanceIdentifier: value.InstanceIdentifier, DisplayName: value.DisplayName,
		Environment: value.Environment, Status: string(value.Status),
	}
}
