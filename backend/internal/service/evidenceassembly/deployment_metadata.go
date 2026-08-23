package evidenceassembly

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type deploymentMetadataPayload struct {
	ProjectName        string `json:"project_name"`
	MonitoringStatus   string `json:"monitoring_status"`
	HealthStatus       string `json:"health_status"`
	RepositoryFullName string `json:"repository_full_name"`
	DefaultBranch      string `json:"default_branch"`
	ApplicationName    string `json:"application_display_name"`
	Environment        string `json:"application_environment"`
	ApplicationStatus  string `json:"application_status"`
}

// deploymentMetadataEvidence snapshots non-secret Project fields already
// safe to expose over REST (ProjectSummaryDTO returns the same data) —
// never credentials.
func deploymentMetadataEvidence(id, investigationID uuid.UUID, project domain.Project, now time.Time) (*domain.Evidence, error) {
	payload := deploymentMetadataPayload{
		ProjectName:        project.Name,
		MonitoringStatus:   string(project.MonitoringStatus),
		HealthStatus:       string(project.HealthStatus),
		RepositoryFullName: project.GitHubRepository.FullName,
		DefaultBranch:      project.GitHubRepository.DefaultBranch,
		ApplicationName:    project.DokployApplication.DisplayName,
		Environment:        project.DokployApplication.Environment,
		ApplicationStatus:  string(project.DokployApplication.Status),
	}
	content, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	summary := fmt.Sprintf(
		"Project %s: repositorio %s@%s, aplicación Dokploy %s (%s) en estado %s.",
		payload.ProjectName, payload.RepositoryFullName, payload.DefaultBranch,
		payload.ApplicationName, payload.Environment, payload.ApplicationStatus,
	)
	return domain.NewEvidence(id, investigationID, domain.EvidenceDeploymentMetadata, summary, string(content), now)
}
