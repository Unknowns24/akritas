package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func MonitoringConfigurationToDomain(value *projectdto.MonitoringConfigurationRequestDTO) (domain.MonitoringConfiguration, error) {
	if value == nil || value.Enabled == nil || value.ErrorPatterns == nil || value.IgnoredPatterns == nil ||
		value.GroupingWindow == nil || value.ContextBefore == nil || value.ContextAfter == nil {
		return domain.MonitoringConfiguration{}, domain.ErrInvalidMonitoringConfiguration
	}
	window, err := parseFixedISODuration(*value.GroupingWindow)
	if err != nil {
		return domain.MonitoringConfiguration{}, domain.ErrInvalidMonitoringConfiguration.Wrap(err)
	}
	return domain.NewMonitoringConfiguration(*value.Enabled, *value.ErrorPatterns, *value.IgnoredPatterns, window, *value.ContextBefore, *value.ContextAfter)
}
