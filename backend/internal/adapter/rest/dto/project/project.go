package projectdto

import inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"

type ProjectDTO struct {
	ProjectSummaryDTO
	MonitoringConfiguration MonitoringConfigurationDTO `json:"monitoring_configuration"`
	BuiltInDetectionRules   []BuiltInDetectionRuleDTO  `json:"built_in_detection_rules"`
}

func FromProject(result *inproject.Result) ProjectDTO {
	summary := FromProjectSummary(*result.Project)
	rules := make([]BuiltInDetectionRuleDTO, 0, len(result.BuiltInDetectionRules))
	for _, rule := range result.BuiltInDetectionRules {
		rules = append(rules, BuiltInDetectionRuleDTO{Code: string(rule.Code), DisplayName: rule.DisplayName, Enabled: rule.Enabled})
	}
	return ProjectDTO{
		ProjectSummaryDTO:       summary,
		MonitoringConfiguration: FromMonitoring(result.Project.MonitoringConfiguration),
		BuiltInDetectionRules:   rules,
	}
}
