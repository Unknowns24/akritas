package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func CreateProjectToCommand(value projectdto.CreateProjectRequestDTO) (portsin.CreateProjectCommand, error) {
	monitoring, err := MonitoringConfigurationToDomain(value.MonitoringConfiguration)
	if err != nil {
		return portsin.CreateProjectCommand{}, err
	}
	return portsin.CreateProjectCommand{Name: value.Name, Description: value.Description, GitHubAccountID: value.GitHubAccountID, RepositoryIdentifier: value.RepositoryIdentifier, DefaultBranch: value.DefaultBranch, DokployServerID: value.DokployServerID, ApplicationIdentifier: value.ApplicationIdentifier, MonitoringConfiguration: monitoring}, nil
}
