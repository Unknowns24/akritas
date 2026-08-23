package domain

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IncidentPhase string

const (
	IncidentPhaseDetected        IncidentPhase = "detected"
	IncidentPhaseInvestigating   IncidentPhase = "investigating"
	IncidentPhasePublishingIssue IncidentPhase = "publishing_issue"
	IncidentPhaseRemediating     IncidentPhase = "remediating"
	IncidentPhaseCompleted       IncidentPhase = "completed"
	IncidentPhaseFailed          IncidentPhase = "failed"
)

func (p IncidentPhase) Validate() error {
	switch p {
	case IncidentPhaseDetected, IncidentPhaseInvestigating, IncidentPhasePublishingIssue, IncidentPhaseRemediating, IncidentPhaseCompleted, IncidentPhaseFailed:
		return nil
	default:
		return ErrInvalidIncidentPhase.Wrap(validationCause("incident phase"))
	}
}

func (p IncidentPhase) IsTerminal() bool {
	return p == IncidentPhaseCompleted || p == IncidentPhaseFailed
}

type TerminalOutcome string

const (
	TerminalOutcomePullRequestCreated     TerminalOutcome = "pull_request_created"
	TerminalOutcomeRequiresHuman          TerminalOutcome = "requires_human"
	TerminalOutcomeRemediationFailed      TerminalOutcome = "remediation_failed"
	TerminalOutcomeInvestigationFailed    TerminalOutcome = "investigation_failed"
	TerminalOutcomeIssuePublicationFailed TerminalOutcome = "issue_publication_failed"
)

func (o TerminalOutcome) Validate() error {
	switch o {
	case TerminalOutcomePullRequestCreated, TerminalOutcomeRequiresHuman, TerminalOutcomeRemediationFailed,
		TerminalOutcomeInvestigationFailed, TerminalOutcomeIssuePublicationFailed:
		return nil
	default:
		return ErrInvalidTerminalOutcome.Wrap(validationCause("terminal outcome"))
	}
}

type RootCauseStatus string

const (
	RootCauseIdentified RootCauseStatus = "identified"
	RootCauseSuspected  RootCauseStatus = "suspected"
	RootCauseUnknown    RootCauseStatus = "unknown"
)

func (s RootCauseStatus) Validate() error {
	switch s {
	case RootCauseIdentified, RootCauseSuspected, RootCauseUnknown:
		return nil
	default:
		return ErrInvalidRootCauseStatus.Wrap(validationCause("root cause status"))
	}
}

type ResolutionStatus string

const (
	ResolutionFixable       ResolutionStatus = "fixable"
	ResolutionRequiresHuman ResolutionStatus = "requires_human"
)

func (s ResolutionStatus) Validate() error {
	switch s {
	case ResolutionFixable, ResolutionRequiresHuman:
		return nil
	default:
		return ErrInvalidResolutionStatus.Wrap(validationCause("resolution status"))
	}
}

type Incident struct {
	ID                   uuid.UUID             `gorm:"column:id;type:uuid;primaryKey"`
	Key                  string                `gorm:"column:key"`
	ProjectID            uuid.UUID             `gorm:"column:project_id;type:uuid"`
	Fingerprint          string                `gorm:"column:fingerprint"`
	Severity             Severity              `gorm:"column:severity"`
	Phase                IncidentPhase         `gorm:"column:phase"`
	TerminalOutcome      *TerminalOutcome      `gorm:"column:terminal_outcome"`
	FirstSeenAt          time.Time             `gorm:"column:first_seen_at"`
	LastSeenAt           time.Time             `gorm:"column:last_seen_at"`
	OccurrenceCount      int64                 `gorm:"column:occurrence_count"`
	Title                string                `gorm:"column:title"`
	Summary              string                `gorm:"column:summary"`
	RootCauseStatus      *RootCauseStatus      `gorm:"column:root_cause_status"`
	ResolutionStatus     *ResolutionStatus     `gorm:"column:resolution_status"`
	Confidence           *float64              `gorm:"column:confidence"`
	GitHubIssueReference *GitHubIssueReference `gorm:"-"`
	PullRequestReference *PullRequestReference `gorm:"serializer:json;type:jsonb;column:pull_request_reference"`
	Project              *ProjectReference     `gorm:"-"`
}

type ProjectReference struct {
	ID   uuid.UUID
	Name string
}

func (i *Incident) PromoteSeverity(severity Severity) {
	if i != nil && severityRank(severity) > severityRank(i.Severity) {
		i.Severity = severity
	}
}

func severityRank(severity Severity) int {
	switch severity {
	case SeverityCritical:
		return 4
	case SeverityError:
		return 3
	case SeverityWarning:
		return 2
	case SeverityInfo:
		return 1
	default:
		return 0
	}
}

var incidentKeyPattern = regexp.MustCompile(`^AKR-[0-9]+$`)

func NewIncident(id uuid.UUID, key string, projectID uuid.UUID, fingerprint string, severity Severity, title string, occurredAt time.Time) (*Incident, error) {
	incident := &Incident{
		ID: id, Key: strings.TrimSpace(key), ProjectID: projectID, Fingerprint: strings.TrimSpace(fingerprint),
		Severity: severity, Phase: IncidentPhaseDetected, FirstSeenAt: occurredAt, LastSeenAt: occurredAt,
		OccurrenceCount: 1, Title: strings.TrimSpace(title),
	}
	if err := incident.Validate(); err != nil {
		return nil, err
	}
	return incident, nil
}

func (i Incident) Validate() error {
	invalid := i.ID == uuid.Nil || !incidentKeyPattern.MatchString(i.Key) || i.ProjectID == uuid.Nil ||
		!nonBlank(i.Fingerprint) || i.Severity.Validate() != nil || i.Phase.Validate() != nil ||
		!validTime(i.FirstSeenAt) || i.LastSeenAt.Before(i.FirstSeenAt) || i.OccurrenceCount < 1 ||
		!nonBlank(i.Title) || len(i.Title) > 500 || len(i.Summary) > 5000
	if invalid {
		return ErrInvalidIncident.Wrap(validationCause("incident"))
	}
	if i.Phase.IsTerminal() != (i.TerminalOutcome != nil) {
		return ErrInvalidIncident.Wrap(validationCause("terminal state"))
	}
	if i.TerminalOutcome != nil && i.TerminalOutcome.Validate() != nil {
		return ErrInvalidIncident.Wrap(validationCause("terminal outcome"))
	}
	if (i.RootCauseStatus == nil) != (i.ResolutionStatus == nil) || (i.ResolutionStatus == nil) != (i.Confidence == nil) {
		return ErrInvalidIncident.Wrap(validationCause("classification"))
	}
	if i.RootCauseStatus != nil && (i.RootCauseStatus.Validate() != nil || i.ResolutionStatus.Validate() != nil || !validConfidence(*i.Confidence)) {
		return ErrInvalidIncident.Wrap(validationCause("classification"))
	}
	if i.GitHubIssueReference != nil && i.GitHubIssueReference.Validate() != nil {
		return ErrInvalidIncident.Wrap(validationCause("issue reference"))
	}
	if i.PullRequestReference != nil && i.PullRequestReference.Validate() != nil {
		return ErrInvalidIncident.Wrap(validationCause("pull request reference"))
	}
	return nil
}

func validConfidence(confidence float64) bool {
	return confidence >= 0 && confidence <= 1
}

func (i Incident) CanGroup(projectID uuid.UUID, fingerprint string, occurredAt time.Time, groupingWindow time.Duration) bool {
	return groupingWindow > 0 && !i.Phase.IsTerminal() && projectID == i.ProjectID && fingerprint == i.Fingerprint &&
		!occurredAt.Before(i.LastSeenAt) && occurredAt.Sub(i.LastSeenAt) <= groupingWindow
}

func (i *Incident) RecordOccurrence(projectID uuid.UUID, fingerprint string, occurredAt time.Time, groupingWindow time.Duration) error {
	if i == nil || !i.CanGroup(projectID, fingerprint, occurredAt, groupingWindow) {
		return ErrIncidentNotGroupable.Wrap(validationCause("incident grouping"))
	}
	i.OccurrenceCount++
	i.LastSeenAt = occurredAt
	return nil
}

func (i *Incident) StartInvestigation() error {
	if i == nil || (i.Phase != IncidentPhaseDetected && i.Phase != IncidentPhaseFailed) {
		return ErrIncidentTransition.Wrap(validationCause("start investigation"))
	}
	i.Phase = IncidentPhaseInvestigating
	i.TerminalOutcome = nil
	return nil
}

func (i *Incident) StartIssuePublication(rootCause RootCauseStatus, resolution ResolutionStatus, confidence float64, summary string) error {
	if i == nil || i.Phase != IncidentPhaseInvestigating || rootCause.Validate() != nil || resolution.Validate() != nil || !validConfidence(confidence) || !nonBlank(summary) {
		return ErrIncidentTransition.Wrap(validationCause("start issue publication"))
	}
	i.RootCauseStatus = &rootCause
	i.ResolutionStatus = &resolution
	i.Confidence = &confidence
	i.Summary = strings.TrimSpace(summary)
	i.Phase = IncidentPhasePublishingIssue
	return nil
}

func (i *Incident) AttachGitHubIssue(reference GitHubIssueReference) error {
	if i == nil || i.Phase != IncidentPhasePublishingIssue || i.GitHubIssueReference != nil || reference.Validate() != nil || reference.IncidentID != i.ID {
		return ErrIncidentTransition.Wrap(validationCause("attach issue"))
	}
	i.GitHubIssueReference = &reference
	return nil
}

func (i *Incident) StartRemediation() error {
	if i == nil || i.Phase != IncidentPhasePublishingIssue || i.GitHubIssueReference == nil ||
		i.ResolutionStatus == nil || *i.ResolutionStatus != ResolutionFixable {
		return ErrIncidentTransition.Wrap(validationCause("start remediation"))
	}
	i.Phase = IncidentPhaseRemediating
	return nil
}

func (i *Incident) CompleteRequiresHuman() error {
	if i == nil || i.Phase != IncidentPhasePublishingIssue || i.GitHubIssueReference == nil ||
		i.ResolutionStatus == nil || *i.ResolutionStatus != ResolutionRequiresHuman {
		return ErrIncidentTransition.Wrap(validationCause("complete requires human"))
	}
	return i.complete(TerminalOutcomeRequiresHuman)
}

func (i *Incident) CompleteRemediationFailed() error {
	if i == nil || i.Phase != IncidentPhaseRemediating || i.GitHubIssueReference == nil {
		return ErrIncidentTransition.Wrap(validationCause("complete remediation failed"))
	}
	return i.complete(TerminalOutcomeRemediationFailed)
}

func (i *Incident) AttachPullRequest(reference PullRequestReference) error {
	if i == nil || i.Phase != IncidentPhaseRemediating || i.PullRequestReference != nil || reference.Validate() != nil {
		return ErrIncidentTransition.Wrap(validationCause("attach pull request"))
	}
	i.PullRequestReference = &reference
	return nil
}

func (i *Incident) CompleteWithPullRequest() error {
	if i == nil || i.Phase != IncidentPhaseRemediating || i.GitHubIssueReference == nil || i.PullRequestReference == nil {
		return ErrIncidentTransition.Wrap(validationCause("complete pull request"))
	}
	return i.complete(TerminalOutcomePullRequestCreated)
}

func (i *Incident) FailInvestigation() error {
	if i == nil || i.Phase != IncidentPhaseInvestigating {
		return ErrIncidentTransition.Wrap(validationCause("fail investigation"))
	}
	return i.fail(TerminalOutcomeInvestigationFailed)
}

func (i *Incident) FailIssuePublication() error {
	if i == nil || i.Phase != IncidentPhasePublishingIssue {
		return ErrIncidentTransition.Wrap(validationCause("fail issue publication"))
	}
	return i.fail(TerminalOutcomeIssuePublicationFailed)
}

func (i *Incident) complete(outcome TerminalOutcome) error {
	i.Phase = IncidentPhaseCompleted
	i.TerminalOutcome = &outcome
	return nil
}

func (i *Incident) fail(outcome TerminalOutcome) error {
	i.Phase = IncidentPhaseFailed
	i.TerminalOutcome = &outcome
	return nil
}
