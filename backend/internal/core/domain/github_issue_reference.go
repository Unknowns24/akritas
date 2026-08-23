package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type GitHubIssueReference struct {
	IncidentID      uuid.UUID `gorm:"column:incident_id;type:uuid"`
	InvestigationID uuid.UUID `gorm:"column:investigation_id;type:uuid;primaryKey"`
	Number          int       `gorm:"column:issue_number"`
	URL             string    `gorm:"column:issue_url"`
	Repository      string    `gorm:"column:repository"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func NewGitHubIssueReference(incidentID, investigationID uuid.UUID, number int, issueURL, repository string, createdAt time.Time) (GitHubIssueReference, error) {
	reference := GitHubIssueReference{
		IncidentID: incidentID, InvestigationID: investigationID, Number: number,
		URL: strings.TrimSpace(issueURL), Repository: strings.TrimSpace(repository), CreatedAt: createdAt,
	}
	if err := reference.Validate(); err != nil {
		return GitHubIssueReference{}, err
	}
	return reference, nil
}

func (r GitHubIssueReference) Validate() error {
	if r.IncidentID == uuid.Nil || r.InvestigationID == uuid.Nil || r.Number < 1 || !validHTTPURL(r.URL) || !nonBlank(r.Repository) || !validTime(r.CreatedAt) {
		return ErrInvalidGitHubIssueReference.Wrap(validationCause("GitHub issue reference"))
	}
	return nil
}
