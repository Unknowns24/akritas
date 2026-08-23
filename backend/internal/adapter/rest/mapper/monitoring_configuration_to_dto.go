package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func MonitoringConfigurationToDTO(value domain.MonitoringConfiguration) projectdto.MonitoringConfigurationDTO {
	errorsOut, ignoredOut := value.ErrorPatterns, value.IgnoredPatterns
	if errorsOut == nil {
		errorsOut = []string{}
	}
	if ignoredOut == nil {
		ignoredOut = []string{}
	}
	return projectdto.MonitoringConfigurationDTO{
		Enabled:         value.Enabled,
		ErrorPatterns:   errorsOut,
		IgnoredPatterns: ignoredOut,
		GroupingWindow:  formatFixedISODuration(value.GroupingWindow),
		ContextBefore:   value.ContextBefore,
		ContextAfter:    value.ContextAfter,
	}
}
