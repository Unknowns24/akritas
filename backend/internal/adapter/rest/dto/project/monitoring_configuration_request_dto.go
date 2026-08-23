package project

type MonitoringConfigurationRequestDTO struct {
	Enabled         *bool     `json:"enabled"`
	ErrorPatterns   *[]string `json:"error_patterns"`
	IgnoredPatterns *[]string `json:"ignored_patterns"`
	GroupingWindow  *string   `json:"grouping_window"`
	ContextBefore   *int      `json:"context_before"`
	ContextAfter    *int      `json:"context_after"`
}
