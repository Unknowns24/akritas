package project

import "github.com/google/uuid"

type CreateProjectRequestDTO struct {
	Name                    string                             `json:"name"`
	Description             string                             `json:"description"`
	GitHubAccountID         uuid.UUID                          `json:"github_account_id"`
	RepositoryIdentifier    string                             `json:"repository_identifier"`
	DefaultBranch           string                             `json:"default_branch"`
	DokployServerID         uuid.UUID                          `json:"dokploy_server_id"`
	ApplicationIdentifier   string                             `json:"application_identifier"`
	MonitoringConfiguration *MonitoringConfigurationRequestDTO `json:"monitoring_configuration"`
}
