package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestIncident(t *testing.T) *Incident {
	t.Helper()
	incident, err := NewIncident(
		uuid.New(),
		"AKR-1",
		uuid.New(),
		"sha256:0123456789abcdef",
		SeverityError,
		"database unavailable",
		time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	return incident
}

func TestIncidentGroupingUsesLastSeenAndInclusiveBoundary(t *testing.T) {
	t.Parallel()

	incident := newTestIncident(t)
	window := 30 * time.Minute
	boundary := incident.LastSeenAt.Add(window)
	if !incident.CanGroup(incident.ProjectID, incident.Fingerprint, boundary, window) {
		t.Fatal("exact grouping boundary should be accepted")
	}
	if err := incident.RecordOccurrence(incident.ProjectID, incident.Fingerprint, boundary, window); err != nil {
		t.Fatal(err)
	}
	if incident.OccurrenceCount != 2 || !incident.LastSeenAt.Equal(boundary) {
		t.Fatalf("occurrence not recorded: %+v", incident)
	}
	if incident.CanGroup(incident.ProjectID, incident.Fingerprint, boundary.Add(window+time.Nanosecond), window) {
		t.Fatal("occurrence outside window should not group")
	}
	if err := incident.RecordOccurrence(uuid.New(), incident.Fingerprint, boundary, window); !errors.Is(err, ErrIncidentNotGroupable) {
		t.Fatalf("expected not groupable error, got %v", err)
	}
}

func TestIncidentRequiresIssueAndValidatedRemediationBeforePR(t *testing.T) {
	t.Parallel()

	incident := newTestIncident(t)
	if err := incident.StartInvestigation(); err != nil {
		t.Fatal(err)
	}
	if err := incident.StartIssuePublication(RootCauseIdentified, ResolutionFixable, 0.9, "root cause"); err != nil {
		t.Fatal(err)
	}
	if err := incident.StartRemediation(); !errors.Is(err, ErrIncidentTransition) {
		t.Fatalf("remediation without issue must fail, got %v", err)
	}
	issue, err := NewGitHubIssueReference(incident.ID, uuid.New(), 7, "https://github.com/Unknowns24/akritas/issues/7", "Unknowns24/akritas", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if err := incident.AttachGitHubIssue(issue); err != nil {
		t.Fatal(err)
	}
	if err := incident.StartRemediation(); err != nil {
		t.Fatal(err)
	}
	if err := incident.CompleteWithPullRequest(); !errors.Is(err, ErrIncidentTransition) {
		t.Fatalf("completion without PR must fail, got %v", err)
	}
	pr, err := NewPullRequestReference(8, "https://github.com/Unknowns24/akritas/pull/8", "Unknowns24/akritas", "codex/fix", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if err := incident.AttachPullRequest(pr); err != nil {
		t.Fatal(err)
	}
	if err := incident.CompleteWithPullRequest(); err != nil {
		t.Fatal(err)
	}
	if incident.Phase != IncidentPhaseCompleted || incident.TerminalOutcome == nil || *incident.TerminalOutcome != TerminalOutcomePullRequestCreated {
		t.Fatalf("unexpected terminal state: %+v", incident)
	}
	if incident.CanGroup(incident.ProjectID, incident.Fingerprint, incident.LastSeenAt, time.Hour) {
		t.Fatal("terminal incident must not group occurrences")
	}
}

func TestIncidentRequiresHumanAndRetryPaths(t *testing.T) {
	t.Parallel()

	incident := newTestIncident(t)
	if err := incident.StartInvestigation(); err != nil {
		t.Fatal(err)
	}
	if err := incident.FailInvestigation(); err != nil {
		t.Fatal(err)
	}
	if err := incident.StartInvestigation(); err != nil {
		t.Fatalf("failed investigation should be retryable: %v", err)
	}
	if incident.TerminalOutcome != nil {
		t.Fatal("retry must clear terminal outcome")
	}
	if err := incident.StartIssuePublication(RootCauseUnknown, ResolutionRequiresHuman, 0, "insufficient evidence"); err != nil {
		t.Fatal(err)
	}
	issue, _ := NewGitHubIssueReference(incident.ID, uuid.New(), 1, "https://example.com/issues/1", "owner/repo", time.Now().UTC())
	if err := incident.AttachGitHubIssue(issue); err != nil {
		t.Fatal(err)
	}
	if err := incident.CompleteRequiresHuman(); err != nil {
		t.Fatal(err)
	}
	if incident.TerminalOutcome == nil || *incident.TerminalOutcome != TerminalOutcomeRequiresHuman {
		t.Fatalf("unexpected outcome: %+v", incident)
	}
}

func TestIncidentFailureOutcomes(t *testing.T) {
	t.Parallel()

	issueFailure := newTestIncident(t)
	_ = issueFailure.StartInvestigation()
	_ = issueFailure.StartIssuePublication(RootCauseUnknown, ResolutionRequiresHuman, 0, "unknown")
	if err := issueFailure.FailIssuePublication(); err != nil {
		t.Fatal(err)
	}
	if issueFailure.Phase != IncidentPhaseFailed || issueFailure.TerminalOutcome == nil || *issueFailure.TerminalOutcome != TerminalOutcomeIssuePublicationFailed {
		t.Fatalf("unexpected issue failure state: %+v", issueFailure)
	}

	remediationFailure := newTestIncident(t)
	_ = remediationFailure.StartInvestigation()
	_ = remediationFailure.StartIssuePublication(RootCauseIdentified, ResolutionFixable, 1, "known")
	issue, _ := NewGitHubIssueReference(remediationFailure.ID, uuid.New(), 2, "https://example.com/issues/2", "owner/repo", time.Now().UTC())
	_ = remediationFailure.AttachGitHubIssue(issue)
	_ = remediationFailure.StartRemediation()
	if err := remediationFailure.CompleteRemediationFailed(); err != nil {
		t.Fatal(err)
	}
	if remediationFailure.TerminalOutcome == nil || *remediationFailure.TerminalOutcome != TerminalOutcomeRemediationFailed {
		t.Fatalf("unexpected remediation failure state: %+v", remediationFailure)
	}
}
