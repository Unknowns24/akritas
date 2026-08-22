package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type ValidationStatus string

const (
	ValidationStatusPending ValidationStatus = "pending"
	ValidationStatusRunning ValidationStatus = "running"
	ValidationStatusPassed  ValidationStatus = "passed"
	ValidationStatusFailed  ValidationStatus = "failed"
)

func (s ValidationStatus) Validate() error {
	switch s {
	case ValidationStatusPending, ValidationStatusRunning, ValidationStatusPassed, ValidationStatusFailed:
		return nil
	default:
		return ErrInvalidValidationStatus.Wrap(validationCause("validation status"))
	}
}

type ValidationType string

const (
	ValidationTypeTest           ValidationType = "test"
	ValidationTypeBuild          ValidationType = "build"
	ValidationTypeStaticAnalysis ValidationType = "static_analysis"
)

func (t ValidationType) Validate() error {
	switch t {
	case ValidationTypeTest, ValidationTypeBuild, ValidationTypeStaticAnalysis:
		return nil
	default:
		return ErrInvalidValidationType.Wrap(validationCause("validation type"))
	}
}

type ValidationResult struct {
	ID             uuid.UUID
	RemediationID  uuid.UUID
	Type           ValidationType
	Name           string
	Status         ValidationStatus
	CreatedAt      time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
	Summary        string
	OutputExcerpt  string
	OutputRedacted bool
}

func NewValidationResult(id, remediationID uuid.UUID, validationType ValidationType, name string, createdAt time.Time) (*ValidationResult, error) {
	result := &ValidationResult{
		ID: id, RemediationID: remediationID, Type: validationType, Name: strings.TrimSpace(name),
		Status: ValidationStatusPending, CreatedAt: createdAt, OutputRedacted: true,
	}
	if err := result.Validate(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r ValidationResult) Validate() error {
	if r.ID == uuid.Nil || r.RemediationID == uuid.Nil || r.Type.Validate() != nil || !nonBlank(r.Name) || len(r.Name) > 255 ||
		r.Status.Validate() != nil || !validTime(r.CreatedAt) || !r.OutputRedacted || len(r.Summary) > 5000 || len(r.OutputExcerpt) > 50000 {
		return ErrInvalidValidationResult.Wrap(validationCause("validation result"))
	}
	switch r.Status {
	case ValidationStatusPending:
		if r.StartedAt != nil || r.FinishedAt != nil {
			return ErrInvalidValidationResult.Wrap(validationCause("pending validation times"))
		}
	case ValidationStatusRunning:
		if r.StartedAt == nil || r.StartedAt.Before(r.CreatedAt) || r.FinishedAt != nil {
			return ErrInvalidValidationResult.Wrap(validationCause("running validation times"))
		}
	case ValidationStatusPassed, ValidationStatusFailed:
		if r.StartedAt == nil || r.FinishedAt == nil || r.StartedAt.Before(r.CreatedAt) || r.FinishedAt.Before(*r.StartedAt) {
			return ErrInvalidValidationResult.Wrap(validationCause("terminal validation times"))
		}
	}
	return nil
}

func (r *ValidationResult) Start(at time.Time) error {
	if r == nil || r.Status != ValidationStatusPending || at.Before(r.CreatedAt) {
		return ErrValidationTransition.Wrap(validationCause("start validation"))
	}
	r.Status = ValidationStatusRunning
	r.StartedAt = &at
	return nil
}

func (r *ValidationResult) Pass(at time.Time, summary, output string) error {
	return r.finish(ValidationStatusPassed, at, summary, output)
}

func (r *ValidationResult) Fail(at time.Time, summary, output string) error {
	return r.finish(ValidationStatusFailed, at, summary, output)
}

func (r *ValidationResult) finish(status ValidationStatus, at time.Time, summary, output string) error {
	if r == nil || r.Status != ValidationStatusRunning || r.StartedAt == nil || at.Before(*r.StartedAt) || !nonBlank(summary) {
		return ErrValidationTransition.Wrap(validationCause("finish validation"))
	}
	r.Status = status
	r.FinishedAt = &at
	r.Summary = strings.TrimSpace(summary)
	r.OutputExcerpt = output
	r.OutputRedacted = true
	return nil
}
