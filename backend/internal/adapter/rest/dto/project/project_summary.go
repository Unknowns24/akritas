package projectdto

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type ProjectSummaryDTO struct {
	ID                 uuid.UUID             `json:"id"`
	Name               string                `json:"name"`
	Description        string                `json:"description,omitempty"`
	MonitoringStatus   string                `json:"monitoring_status"`
	HealthStatus       string                `json:"health_status"`
	GitHubRepository   GitHubRepositoryDTO   `json:"github_repository"`
	DokployApplication DokployApplicationDTO `json:"dokploy_application"`
	LastObservedAt     *time.Time            `json:"last_observed_at,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

func FromProjectSummary(project domain.Project) ProjectSummaryDTO {
	return ProjectSummaryDTO{
		ID: project.ID, Name: project.Name, Description: project.Description,
		MonitoringStatus: string(project.MonitoringStatus), HealthStatus: string(project.HealthStatus),
		GitHubRepository: GitHubRepositoryDTO{
			GitHubAccountID: project.GitHubRepository.GitHubAccountID, RepositoryIdentifier: project.GitHubRepository.RepositoryIdentifier,
			Owner: project.GitHubRepository.Owner, Name: project.GitHubRepository.Name, FullName: project.GitHubRepository.FullName,
			DefaultBranch: project.GitHubRepository.DefaultBranch, Private: project.GitHubRepository.Private,
			HTMLURL: project.GitHubRepository.HTMLURL,
		},
		DokployApplication: DokployApplicationDTO{
			DokployServerID: project.DokployApplication.DokployServerID, ApplicationIdentifier: project.DokployApplication.ApplicationIdentifier,
			InstanceIdentifier: project.DokployApplication.InstanceIdentifier, DisplayName: project.DokployApplication.DisplayName,
			Environment: project.DokployApplication.Environment, Status: string(project.DokployApplication.Status),
		},
		LastObservedAt: project.LastObservedAt, CreatedAt: project.CreatedAt.UTC(), UpdatedAt: project.UpdatedAt.UTC(),
	}
}
