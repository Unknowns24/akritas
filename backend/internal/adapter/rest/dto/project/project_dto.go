package project

type ProjectDTO struct {
	ProjectSummaryDTO
	MonitoringConfiguration MonitoringConfigurationDTO `json:"monitoring_configuration"`
	BuiltInDetectionRules   []BuiltInDetectionRuleDTO  `json:"built_in_detection_rules"`
}
