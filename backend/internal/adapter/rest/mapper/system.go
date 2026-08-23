package mapper

import (
	systemdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/system"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func SystemStatusToDTO(value domain.SystemStatus) systemdto.SystemStatusDTO {
	components := make([]systemdto.ComponentHealthDTO, 0, len(value.Components))
	for _, component := range value.Components {
		components = append(components, systemdto.ComponentHealthDTO{
			Component: component.Component,
			Status:    string(component.Status),
			CheckedAt: component.CheckedAt,
		})
	}
	return systemdto.SystemStatusDTO{
		GitHubAccountCount: value.GitHubAccountCount,
		DokployServerCount: value.DokployServerCount,
		QvacEndpoint:       value.QvacEndpoint,
		Components:         components,
		LastDiagnosticsAt:  value.LastDiagnosticsAt,
	}
}
