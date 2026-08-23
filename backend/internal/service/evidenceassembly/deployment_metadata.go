package evidenceassembly

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
	"github.com/google/uuid"
)

type deploymentMetadataPayload struct {
	ProjectName        string `json:"project_name"`
	MonitoringStatus   string `json:"monitoring_status"`
	HealthStatus       string `json:"health_status"`
	RepositoryFullName string `json:"repository_full_name"`
	DefaultBranch      string `json:"default_branch"`
	SourceType         string `json:"source_type"`
	SourceDisplayName  string `json:"source_display_name"`
	SourceEnvironment  string `json:"source_environment"`
	SourceStatus       string `json:"source_status"`
	SourceServiceName  string `json:"source_service_name,omitempty"`
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
		SourceType:         string(project.DokploySource.Type),
		SourceDisplayName:  project.DokploySource.DisplayName,
		SourceEnvironment:  project.DokploySource.Environment,
		SourceStatus:       string(project.DokploySource.Status),
		SourceServiceName:  project.DokploySource.ServiceName,
	}
	content, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	summary := fmt.Sprintf(
		"Project %s: repositorio %s@%s, fuente Dokploy %s (%s) en estado %s.",
		payload.ProjectName, payload.RepositoryFullName, payload.DefaultBranch,
		payload.SourceDisplayName, payload.SourceEnvironment, payload.SourceStatus,
	)
	return domain.NewEvidence(id, investigationID, domain.EvidenceDeploymentMetadata, evidencesafety.Redact(summary), evidencesafety.Redact(string(content)), now)
}
