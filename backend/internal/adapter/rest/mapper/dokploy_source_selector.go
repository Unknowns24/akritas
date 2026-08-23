package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokploySourceSelectorToDomain(value projectdto.DokploySourceRequestDTO) domain.DokploySourceSelector {
	return domain.DokploySourceSelector{Type: domain.DokploySourceType(value.Type), DokployServerID: value.DokployServerID, ResourceIdentifier: value.ResourceIdentifier, ServiceName: value.ServiceName}.Normalize()
}
