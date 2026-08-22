package projectdto

import (
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/google/uuid"
)

type CreateProjectRequestDTO struct {
	Name                    string                     `json:"name"`
	Description             string                     `json:"description"`
	GitHubAccountID         uuid.UUID                  `json:"github_account_id"`
	RepositoryIdentifier    string                     `json:"repository_identifier"`
	DefaultBranch           string                     `json:"default_branch"`
	DokployServerID         uuid.UUID                  `json:"dokploy_server_id"`
	ApplicationIdentifier   string                     `json:"application_identifier"`
	MonitoringConfiguration MonitoringConfigurationDTO `json:"monitoring_configuration"`
}

func (r CreateProjectRequestDTO) Command() (project.CreateCommand, error) {
	configuration, err := r.MonitoringConfiguration.Domain()
	if err != nil {
		return project.CreateCommand{}, err
	}
	return project.CreateCommand{
		Name: r.Name, Description: r.Description, GitHubAccountID: r.GitHubAccountID,
		RepositoryIdentifier: r.RepositoryIdentifier, DefaultBranch: r.DefaultBranch,
		DokployServerID: r.DokployServerID, ApplicationIdentifier: r.ApplicationIdentifier,
		MonitoringConfiguration: configuration,
	}, nil
}
