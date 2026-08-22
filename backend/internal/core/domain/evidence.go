package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type EvidenceType string

const (
	EvidenceLogExcerpt         EvidenceType = "log_excerpt"
	EvidenceStackTrace         EvidenceType = "stack_trace"
	EvidenceCodeLocation       EvidenceType = "code_location"
	EvidenceCommit             EvidenceType = "commit"
	EvidenceDiff               EvidenceType = "diff"
	EvidenceValidationResult   EvidenceType = "validation_result"
	EvidenceDeploymentMetadata EvidenceType = "deployment_metadata"
)

func (t EvidenceType) Validate() error {
	switch t {
	case EvidenceLogExcerpt, EvidenceStackTrace, EvidenceCodeLocation, EvidenceCommit, EvidenceDiff,
		EvidenceValidationResult, EvidenceDeploymentMetadata:
		return nil
	default:
		return ErrInvalidEvidenceType.Wrap(validationCause("evidence type"))
	}
}

type Evidence struct {
	ID              uuid.UUID
	InvestigationID uuid.UUID
	Type            EvidenceType
	Summary         string
	Content         string
	FilePath        string
	LineStart       *int
	LineEnd         *int
	CommitSHA       string
	Patch           string
	OccurredAt      *time.Time
	Redacted        bool
	CreatedAt       time.Time
}

func NewEvidence(id, investigationID uuid.UUID, evidenceType EvidenceType, summary, content string, createdAt time.Time) (*Evidence, error) {
	evidence := &Evidence{
		ID: id, InvestigationID: investigationID, Type: evidenceType,
		Summary: strings.TrimSpace(summary), Content: content, Redacted: true, CreatedAt: createdAt,
	}
	if err := evidence.Validate(); err != nil {
		return nil, err
	}
	return evidence, nil
}

func (e Evidence) Validate() error {
	if e.ID == uuid.Nil || e.InvestigationID == uuid.Nil || e.Type.Validate() != nil || !nonBlank(e.Summary) ||
		len(e.Summary) > 5000 || len(e.Content) > 50000 || len(e.FilePath) > 4096 || len(e.CommitSHA) > 64 ||
		len(e.Patch) > 100000 || !e.Redacted || !validTime(e.CreatedAt) {
		return ErrInvalidEvidence.Wrap(validationCause("evidence"))
	}
	if (e.LineStart == nil) != (e.LineEnd == nil) || (e.LineStart != nil && (*e.LineStart < 1 || *e.LineEnd < *e.LineStart)) {
		return ErrInvalidEvidence.Wrap(validationCause("evidence lines"))
	}
	return nil
}
