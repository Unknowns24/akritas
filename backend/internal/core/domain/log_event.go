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
	ID                  uuid.UUID            `gorm:"column:id;type:uuid;primaryKey"`
	IncidentID          uuid.UUID            `gorm:"column:incident_id;type:uuid"`
	ProjectID           uuid.UUID            `gorm:"column:project_id;type:uuid"`
	SourceApplicationID string               `gorm:"column:source_application_id"`
	SourceInstanceID    string               `gorm:"column:source_instance_id"`
	OccurrenceKey       string               `gorm:"column:occurrence_key"`
	Timestamp           time.Time            `gorm:"column:timestamp"`
	Severity            Severity             `gorm:"column:severity"`
	Message             string               `gorm:"column:message"`
	Fingerprint         string               `gorm:"column:fingerprint"`
	DetectionRules      []string             `gorm:"serializer:json;type:jsonb;column:detection_rules"`
	ContextBefore       []SanitizedLogRecord `gorm:"serializer:json;type:jsonb;column:context_before"`
	ContextAfter        []SanitizedLogRecord `gorm:"serializer:json;type:jsonb;column:context_after"`
	RawContextRedacted  bool                 `gorm:"column:raw_context_redacted"`
}

func (e *LogEvent) AssociateOccurrence(incidentID uuid.UUID, sourceApplicationID, sourceInstanceID, occurrenceKey string) error {
	if e == nil || incidentID == uuid.Nil || strings.TrimSpace(sourceApplicationID) == "" || strings.TrimSpace(sourceInstanceID) == "" || strings.TrimSpace(occurrenceKey) == "" {
		return ErrInvalidLogEvent.Wrap(validationCause("log event occurrence"))
	}
	e.IncidentID = incidentID
	e.SourceApplicationID = strings.TrimSpace(sourceApplicationID)
	e.SourceInstanceID = strings.TrimSpace(sourceInstanceID)
	e.OccurrenceKey = strings.TrimSpace(occurrenceKey)
	return nil
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
