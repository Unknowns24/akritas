package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityError    Severity = "error"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)

func (s Severity) Validate() error {
	switch s {
	case SeverityCritical, SeverityError, SeverityWarning, SeverityInfo:
		return nil
	default:
		return ErrInvalidSeverity.Wrap(validationCause("severity"))
	}
}

type LogEvent struct {
	ID                 uuid.UUID
	ProjectID          uuid.UUID
	Timestamp          time.Time
	Severity           Severity
	Message            string
	Fingerprint        string
	DetectionRules     []string
	ContextBefore      []SanitizedLogRecord
	ContextAfter       []SanitizedLogRecord
	RawContextRedacted bool
}

func NewLogEvent(
	id, projectID uuid.UUID,
	timestamp time.Time,
	severity Severity,
	message, fingerprint string,
	detectionRules []string,
	contextBefore, contextAfter []SanitizedLogRecord,
) (*LogEvent, error) {
	event := &LogEvent{
		ID: id, ProjectID: projectID, Timestamp: timestamp, Severity: severity,
		Message: strings.TrimSpace(message), Fingerprint: strings.TrimSpace(fingerprint),
		DetectionRules:     cloneStrings(detectionRules),
		ContextBefore:      append([]SanitizedLogRecord(nil), contextBefore...),
		ContextAfter:       append([]SanitizedLogRecord(nil), contextAfter...),
		RawContextRedacted: true,
	}
	if err := event.Validate(); err != nil {
		return nil, err
	}
	return event, nil
}

func (e LogEvent) Validate() error {
	if e.ID == uuid.Nil || e.ProjectID == uuid.Nil || !validTime(e.Timestamp) || e.Severity.Validate() != nil ||
		!nonBlank(e.Message) || len(e.Message) > 20000 || !nonBlank(e.Fingerprint) || len(e.DetectionRules) == 0 ||
		!e.RawContextRedacted {
		return ErrInvalidLogEvent.Wrap(validationCause("log event"))
	}
	for _, rule := range e.DetectionRules {
		if !nonBlank(rule) {
			return ErrInvalidLogEvent.Wrap(validationCause("detection rule"))
		}
	}
	for _, record := range append(append([]SanitizedLogRecord(nil), e.ContextBefore...), e.ContextAfter...) {
		if record.Validate() != nil {
			return ErrInvalidLogEvent.Wrap(validationCause("log context"))
		}
	}
	return nil
}
