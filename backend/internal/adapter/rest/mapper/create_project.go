package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func CreateProjectToCommand(value projectdto.CreateProjectRequestDTO) (portsin.CreateProjectCommand, error) {
	monitoring, err := MonitoringConfigurationToDomain(value.MonitoringConfiguration)
	if err != nil {
		return portsin.CreateProjectCommand{}, err
	}
	ingestion := domain.InitialLogIngestion(value.InitialLogIngestion)
	if ingestion == "" {
		ingestion = domain.InitialLogIngestionFromNow
	}
	if err := ingestion.Validate(); err != nil {
		return portsin.CreateProjectCommand{}, err
	}
	return portsin.CreateProjectCommand{Name: value.Name, Description: value.Description, GitHubAccountID: value.GitHubAccountID, RepositoryIdentifier: value.RepositoryIdentifier, DefaultBranch: value.DefaultBranch, DokployServerID: value.DokployServerID, ApplicationIdentifier: value.ApplicationIdentifier, MonitoringConfiguration: monitoring, InitialLogIngestion: ingestion}, nil
}
