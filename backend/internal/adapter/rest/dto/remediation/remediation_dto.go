package remediation

import (
	"time"

	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/google/uuid"
)

type ValidationSummaryDTO struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
}

type CodeChangeDTO struct {
	FilePath   string `json:"file_path"`
	ChangeType string `json:"change_type"`
	Patch      string `json:"patch"`
	Redacted   bool   `json:"redacted"`
}

type PullRequestReferenceDTO struct {
	Number     int       `json:"number"`
	URL        string    `json:"url"`
	Repository string    `json:"repository"`
	Branch     string    `json:"branch"`
	CreatedAt  time.Time `json:"created_at"`
}

type ValidationResultDTO struct {
	ID             uuid.UUID  `json:"id"`
	RemediationID  uuid.UUID  `json:"remediation_id"`
	Type           string     `json:"type"`
	Name           string     `json:"name"`
	Status         string     `json:"status"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
	Summary        string     `json:"summary,omitempty"`
	OutputExcerpt  string     `json:"output_excerpt,omitempty"`
	OutputRedacted bool       `json:"output_redacted"`
}

type RemediationDTO struct {
	ID                   uuid.UUID                `json:"id"`
	IncidentID           uuid.UUID                `json:"incident_id"`
	Status               string                   `json:"status"`
	BranchName           string                   `json:"branch_name,omitempty"`
	ChangesSummary       string                   `json:"changes_summary,omitempty"`
	Changes              []CodeChangeDTO          `json:"changes"`
	ValidationSummary    ValidationSummaryDTO     `json:"validation_summary"`
	FailureUserMessage   string                   `json:"failure_user_message,omitempty"`
	PullRequestReference *PullRequestReferenceDTO `json:"pull_request_reference,omitempty"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
}

type PullRequestDTO struct {
	ID                uuid.UUID                            `json:"id"`
	Project           incidentdto.ProjectReferenceDTO      `json:"project"`
	IncidentID        uuid.UUID                            `json:"incident_id"`
	IncidentKey       string                               `json:"incident_key"`
	RemediationID     uuid.UUID                            `json:"remediation_id"`
	IssueReference    *incidentdto.GitHubIssueReferenceDTO `json:"issue_reference"`
	Reference         PullRequestReferenceDTO              `json:"reference"`
	Title             string                               `json:"title,omitempty"`
	ChangesSummary    string                               `json:"changes_summary,omitempty"`
	ValidationSummary string                               `json:"validation_summary,omitempty"`
	RiskSummary       string                               `json:"risk_summary,omitempty"`
	CreatedAt         time.Time                            `json:"created_at"`
}
