package investigation

import "time"

type InvestigationDTO struct {
	ID                 string     `json:"id"`
	IncidentID         string     `json:"incident_id"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	StartedAt          *time.Time `json:"started_at,omitempty"`
	FinishedAt         *time.Time `json:"finished_at,omitempty"`
	Summary            string     `json:"summary,omitempty"`
	RootCause          string     `json:"root_cause,omitempty"`
	RootCauseStatus    *string    `json:"root_cause_status,omitempty"`
	ResolutionStatus   *string    `json:"resolution_status,omitempty"`
	Confidence         *float64   `json:"confidence,omitempty"`
	Hypotheses         []string   `json:"hypotheses"`
	RelevantFiles      []string   `json:"relevant_files"`
	RelevantCommits    []string   `json:"relevant_commits"`
	RecommendedActions []string   `json:"recommended_actions"`
	EvidenceIDs        []string   `json:"evidence_ids"`
	EvidenceCount      int        `json:"evidence_count,omitempty"`
	FailureUserMessage string     `json:"failure_user_message,omitempty"`
}
