package project

import "github.com/google/uuid"

type CreateProjectRequestDTO struct {
	Name                    string                             `json:"name"`
	Description             string                             `json:"description"`
	GitHubAccountID         uuid.UUID                          `json:"github_account_id"`
	RepositoryIdentifier    string                             `json:"repository_identifier"`
	DefaultBranch           string                             `json:"default_branch"`
	DokploySource           DokploySourceRequestDTO            `json:"dokploy_source"`
	MonitoringConfiguration *MonitoringConfigurationRequestDTO `json:"monitoring_configuration"`
	InitialLogIngestion     string                             `json:"initial_log_ingestion,omitempty"`
}
