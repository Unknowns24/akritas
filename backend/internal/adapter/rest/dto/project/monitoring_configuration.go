package projectdto

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/utils"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type MonitoringConfigurationDTO struct {
	Enabled         bool     `json:"enabled"`
	ErrorPatterns   []string `json:"error_patterns"`
	IgnoredPatterns []string `json:"ignored_patterns"`
	GroupingWindow  string   `json:"grouping_window"`
	ContextBefore   int      `json:"context_before"`
	ContextAfter    int      `json:"context_after"`
}

func (c MonitoringConfigurationDTO) Domain() (domain.MonitoringConfiguration, error) {
	window, err := utils.ParseISODuration(c.GroupingWindow)
	if err != nil {
		return domain.MonitoringConfiguration{}, domain.ErrInvalidMonitoringConfiguration.Wrap(err)
	}
	errorPatterns := c.ErrorPatterns
	if errorPatterns == nil {
		errorPatterns = []string{}
	}
	ignoredPatterns := c.IgnoredPatterns
	if ignoredPatterns == nil {
		ignoredPatterns = []string{}
	}
	return domain.NewMonitoringConfiguration(c.Enabled, errorPatterns, ignoredPatterns, window, c.ContextBefore, c.ContextAfter)
}

func FromMonitoring(configuration domain.MonitoringConfiguration) MonitoringConfigurationDTO {
	errorPatterns := configuration.ErrorPatterns
	if errorPatterns == nil {
		errorPatterns = []string{}
	}
	ignoredPatterns := configuration.IgnoredPatterns
	if ignoredPatterns == nil {
		ignoredPatterns = []string{}
	}
	return MonitoringConfigurationDTO{
		Enabled: configuration.Enabled, ErrorPatterns: errorPatterns,
		IgnoredPatterns: ignoredPatterns, GroupingWindow: utils.FormatISODuration(configuration.GroupingWindow),
		ContextBefore: configuration.ContextBefore, ContextAfter: configuration.ContextAfter,
	}
}
