package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func ProjectResultToDTO(value *portsin.ProjectResult) projectdto.ProjectDTO {
	rules := make([]projectdto.BuiltInDetectionRuleDTO, 0, len(value.BuiltInDetectionRules))
	for _, rule := range value.BuiltInDetectionRules {
		rules = append(rules, projectdto.BuiltInDetectionRuleDTO{
			Code:        string(rule.Code),
			DisplayName: rule.DisplayName,
			Enabled:     rule.Enabled,
		})
	}
	return projectdto.ProjectDTO{
		ProjectSummaryDTO:       ProjectSummaryToDTO(*value.Project),
		MonitoringConfiguration: MonitoringConfigurationToDTO(value.Project.MonitoringConfiguration),
		BuiltInDetectionRules:   rules,
	}
}
