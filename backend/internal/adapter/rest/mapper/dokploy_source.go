package mapper

import (
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokploySourceToDTO(value domain.DokploySource) dokploydto.DokploySourceDTO {
	return dokploydto.DokploySourceDTO{
		Type: string(value.Type), DokployServerID: value.DokployServerID.String(), ResourceIdentifier: value.ResourceIdentifier,
		ServiceName: value.ServiceName, InstanceIdentifier: value.InstanceIdentifier, DisplayName: value.DisplayName,
		Environment: value.Environment, Status: string(value.Status), RuntimeType: string(value.RuntimeType),
	}
}
