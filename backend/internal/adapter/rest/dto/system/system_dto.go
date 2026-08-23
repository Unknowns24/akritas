package system

import "time"

type HealthDTO struct {
	Status string `json:"status"`
}

type ReadinessDTO struct {
	Status string `json:"status"`
}

type ComponentHealthDTO struct {
	Component string     `json:"component"`
	Status    string     `json:"status"`
	CheckedAt *time.Time `json:"checked_at,omitempty"`
}

type SystemStatusDTO struct {
	GitHubAccountCount int                  `json:"github_account_count"`
	DokployServerCount int                  `json:"dokploy_server_count"`
	QvacEndpoint       string               `json:"qvac_endpoint,omitempty"`
	Components         []ComponentHealthDTO `json:"components"`
	LastDiagnosticsAt  *time.Time           `json:"last_diagnostics_at,omitempty"`
}
