package evidence

import "time"

type EvidenceDTO struct {
	ID              string     `json:"id"`
	InvestigationID string     `json:"investigation_id"`
	Type            string     `json:"type"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content,omitempty"`
	FilePath        string     `json:"file_path,omitempty"`
	LineStart       *int       `json:"line_start,omitempty"`
	LineEnd         *int       `json:"line_end,omitempty"`
	CommitSHA       string     `json:"commit_sha,omitempty"`
	Patch           string     `json:"patch,omitempty"`
	OccurredAt      *time.Time `json:"occurred_at,omitempty"`
	Redacted        bool       `json:"redacted"`
	CreatedAt       time.Time  `json:"created_at"`
}
