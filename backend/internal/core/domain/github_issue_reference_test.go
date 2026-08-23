package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGitHubIssueReferenceRequiresIncidentAndInvestigationLinks(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	incidentID, investigationID := uuid.New(), uuid.New()
	reference, err := NewGitHubIssueReference(incidentID, investigationID, 42,
		"https://github.com/Unknowns24/akritas/issues/42", "Unknowns24/akritas", now)
	if err != nil {
		t.Fatal(err)
	}
	if reference.IncidentID != incidentID || reference.InvestigationID != investigationID || reference.Number != 42 || reference.CreatedAt != now {
		t.Fatalf("reference=%+v", reference)
	}
	for _, invalid := range []GitHubIssueReference{
		{InvestigationID: investigationID, Number: 1, URL: reference.URL, Repository: reference.Repository, CreatedAt: now},
		{IncidentID: incidentID, Number: 1, URL: reference.URL, Repository: reference.Repository, CreatedAt: now},
	} {
		if err := invalid.Validate(); !errors.Is(err, ErrInvalidGitHubIssueReference) {
			t.Fatalf("expected invalid linked reference, got %v", err)
		}
	}
}
