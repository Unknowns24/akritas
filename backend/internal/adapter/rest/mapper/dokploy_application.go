package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokployApplicationToDTO(value domain.DokployApplication) dto.DokployApplicationDTO {
	return dto.DokployApplicationDTO{
		DokployServerID: value.DokployServerID.String(), ApplicationIdentifier: value.ApplicationIdentifier,
		InstanceIdentifier: value.InstanceIdentifier, DisplayName: value.DisplayName,
		Environment: value.Environment, Status: string(value.Status),
	}
}
