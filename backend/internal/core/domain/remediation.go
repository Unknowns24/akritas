package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type RemediationStatus string

const (
	RemediationStatusPlanned            RemediationStatus = "planned"
	RemediationStatusInProgress         RemediationStatus = "in_progress"
	RemediationStatusValidated          RemediationStatus = "validated"
	RemediationStatusFailed             RemediationStatus = "failed"
	RemediationStatusPullRequestCreated RemediationStatus = "pull_request_created"
)

func (s RemediationStatus) Validate() error {
	switch s {
	case RemediationStatusPlanned, RemediationStatusInProgress, RemediationStatusValidated,
		RemediationStatusFailed, RemediationStatusPullRequestCreated:
		return nil
	default:
		return ErrInvalidRemediationStatus.Wrap(validationCause("remediation status"))
	}
}

type Remediation struct {
	ID                   uuid.UUID
	IncidentID           uuid.UUID
	InvestigationID      uuid.UUID
	Status               RemediationStatus
	BranchName           string
	ChangesSummary       string
	Changes              []CodeChange
	ValidationResults    []ValidationResult
	FailureUserMessage   string
	PullRequestReference *PullRequestReference
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewRemediation(id, incidentID uuid.UUID, createdAt time.Time) (*Remediation, error) {
	remediation := &Remediation{
		ID: id, IncidentID: incidentID, Status: RemediationStatusPlanned,
		Changes: []CodeChange{}, ValidationResults: []ValidationResult{}, CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	if err := remediation.Validate(); err != nil {
		return nil, err
	}
	return remediation, nil
}

func (r Remediation) Validate() error {
	if r.ID == uuid.Nil || r.IncidentID == uuid.Nil || r.Status.Validate() != nil || !validTime(r.CreatedAt) || r.UpdatedAt.Before(r.CreatedAt) {
		return ErrInvalidRemediation.Wrap(validationCause("remediation"))
	}
	if r.Status != RemediationStatusPlanned && !nonBlank(r.BranchName) {
		return ErrInvalidRemediation.Wrap(validationCause("remediation branch"))
	}
	if r.Status == RemediationStatusFailed && !nonBlank(r.FailureUserMessage) {
		return ErrInvalidRemediation.Wrap(validationCause("remediation failure"))
	}
	if (r.Status == RemediationStatusPullRequestCreated) != (r.PullRequestReference != nil) {
		return ErrInvalidRemediation.Wrap(validationCause("remediation pull request"))
	}
	for _, change := range r.Changes {
		if change.Validate() != nil {
			return ErrInvalidRemediation.Wrap(validationCause("remediation change"))
		}
	}
	for _, result := range r.ValidationResults {
		if result.RemediationID != r.ID || result.Validate() != nil {
			return ErrInvalidRemediation.Wrap(validationCause("remediation validation"))
		}
	}
	return nil
}

func (r *Remediation) AttachInvestigation(investigationID uuid.UUID, at time.Time) error {
	if r == nil || r.InvestigationID != uuid.Nil || investigationID == uuid.Nil || at.Before(r.UpdatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("attach remediation investigation"))
	}
	r.InvestigationID = investigationID
	r.UpdatedAt = at
	return nil
}

func (r *Remediation) Start(branchName string, at time.Time) error {
	if r == nil || r.Status != RemediationStatusPlanned || !nonBlank(branchName) || at.Before(r.CreatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("start remediation"))
	}
	r.Status = RemediationStatusInProgress
	r.BranchName = strings.TrimSpace(branchName)
	r.UpdatedAt = at
	return nil
}

func (r *Remediation) AddChange(change CodeChange, at time.Time) error {
	if r == nil || r.Status != RemediationStatusInProgress || change.Validate() != nil || at.Before(r.UpdatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("add code change"))
	}
	r.Changes = append(r.Changes, change)
	r.UpdatedAt = at
	return nil
}

func (r *Remediation) AddValidationResult(result ValidationResult, at time.Time) error {
	if r == nil || r.Status != RemediationStatusInProgress || result.RemediationID != r.ID || result.Validate() != nil || at.Before(r.UpdatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("add validation result"))
	}
	r.ValidationResults = append(r.ValidationResults, result)
	r.UpdatedAt = at
	return nil
}

func (r *Remediation) MarkValidated(at time.Time) error {
	if r == nil || r.Status != RemediationStatusInProgress || len(r.Changes) == 0 || len(r.ValidationResults) == 0 || at.Before(r.UpdatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("validate remediation"))
	}
	for _, result := range r.ValidationResults {
		if result.Status != ValidationStatusPassed {
			return ErrRemediationTransition.Wrap(validationCause("validation did not pass"))
		}
	}
	r.Status = RemediationStatusValidated
	r.UpdatedAt = at
	return nil
}

func (r *Remediation) Fail(userMessage string, at time.Time) error {
	if r == nil || (r.Status != RemediationStatusInProgress && r.Status != RemediationStatusValidated) || !nonBlank(userMessage) || at.Before(r.UpdatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("fail remediation"))
	}
	r.Status = RemediationStatusFailed
	r.FailureUserMessage = strings.TrimSpace(userMessage)
	r.UpdatedAt = at
	return nil
}

func (r *Remediation) CreatePullRequest(reference PullRequestReference, at time.Time) error {
	if r == nil || r.Status != RemediationStatusValidated || reference.Validate() != nil || at.Before(r.UpdatedAt) {
		return ErrRemediationTransition.Wrap(validationCause("create pull request"))
	}
	r.Status = RemediationStatusPullRequestCreated
	r.PullRequestReference = &reference
	r.UpdatedAt = at
	return nil
}
