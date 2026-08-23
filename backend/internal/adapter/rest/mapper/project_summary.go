package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func ProjectSummaryToDTO(value domain.Project) projectdto.ProjectSummaryDTO {
	return projectdto.ProjectSummaryDTO{
		ID:               value.ID.String(),
		Name:             value.Name,
		Description:      value.Description,
		MonitoringStatus: string(value.MonitoringStatus),
		HealthStatus:     string(value.HealthStatus),
		GitHubRepository: GitHubRepositoryToDTO(value.GitHubRepository),
		DokploySource:    DokploySourceToDTO(value.DokploySource),
		LastObservedAt:   value.LastObservedAt,
		CreatedAt:        value.CreatedAt,
		UpdatedAt:        value.UpdatedAt,
	}
}
