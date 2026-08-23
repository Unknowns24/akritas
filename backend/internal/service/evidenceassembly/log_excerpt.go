package evidenceassembly

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
	"github.com/google/uuid"
)

const maximumLogEvidenceBytes = 16 << 10

type logRecordPayload struct {
	Timestamp time.Time `json:"timestamp"`
	Stream    string    `json:"stream"`
	Message   string    `json:"message"`
	Redacted  bool      `json:"redacted"`
}

type logExcerptPayload struct {
	Message             string             `json:"message"`
	Timestamp           time.Time          `json:"timestamp"`
	Severity            string             `json:"severity"`
	DetectionRules      []string           `json:"detection_rules"`
	SourceApplicationID string             `json:"source_application_id,omitempty"`
	SourceInstanceID    string             `json:"source_instance_id,omitempty"`
	ContextBefore       []logRecordPayload `json:"context_before"`
	ContextAfter        []logRecordPayload `json:"context_after"`
	Redacted            bool               `json:"redacted"`
}

func logExcerptEvidence(id, investigationID uuid.UUID, event domain.LogEvent, now time.Time) (*domain.Evidence, error) {
	payload := logExcerptPayload{
		Message:             evidencesafety.Redact(event.Message),
		Timestamp:           event.Timestamp,
		Severity:            string(event.Severity),
		DetectionRules:      append([]string(nil), event.DetectionRules...),
		SourceApplicationID: evidencesafety.Redact(event.SourceApplicationID),
		SourceInstanceID:    evidencesafety.Redact(event.SourceInstanceID),
		ContextBefore:       logRecords(event.ContextBefore),
		ContextAfter:        logRecords(event.ContextAfter),
		Redacted:            true,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	evidence, err := domain.NewEvidence(id, investigationID, domain.EvidenceLogExcerpt,
		fmt.Sprintf("Log %s detectado a las %s.", event.Severity, event.Timestamp.UTC().Format(time.RFC3339)),
		evidencesafety.RedactAndLimit(string(raw), maximumLogEvidenceBytes), now)
	if err != nil {
		return nil, err
	}
	evidence.OccurredAt = &event.Timestamp
	return evidence, evidence.Validate()
}

func stackTraceEvidence(id, investigationID uuid.UUID, event domain.LogEvent, now time.Time) (*domain.Evidence, error) {
	var stack strings.Builder
	for _, record := range event.ContextBefore {
		stack.WriteString(record.Message)
		stack.WriteByte('\n')
	}
	stack.WriteString(event.Message)
	for _, record := range event.ContextAfter {
		stack.WriteByte('\n')
		stack.WriteString(record.Message)
	}
	content := evidencesafety.RedactAndLimit(stack.String(), maximumLogEvidenceBytes)
	evidence, err := domain.NewEvidence(id, investigationID, domain.EvidenceStackTrace,
		fmt.Sprintf("Stack trace real detectado a las %s.", event.Timestamp.UTC().Format(time.RFC3339)), content, now)
	if err != nil {
		return nil, err
	}
	evidence.OccurredAt = &event.Timestamp
	return evidence, evidence.Validate()
}

func logRecords(records []domain.SanitizedLogRecord) []logRecordPayload {
	result := make([]logRecordPayload, 0, len(records))
	for _, record := range records {
		result = append(result, logRecordPayload{
			Timestamp: record.Timestamp,
			Stream:    string(record.Stream),
			Message:   evidencesafety.Redact(record.Message),
			Redacted:  true,
		})
	}
	return result
}

func hasStackTrace(event domain.LogEvent) bool {
	return slices.Contains(event.DetectionRules, string(domain.DetectionRuleStackTrace))
}
