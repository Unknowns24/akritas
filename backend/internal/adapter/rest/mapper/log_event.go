package mapper

import (
	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func LogEventToDTO(value domain.LogEvent) incidentdto.LogEventDTO {
	return incidentdto.LogEventDTO{
		ID: value.ID, ProjectID: value.ProjectID, Timestamp: value.Timestamp, Severity: value.Severity,
		Message: value.Message, Fingerprint: value.Fingerprint, DetectionRules: append([]string(nil), value.DetectionRules...),
		ContextBefore: logRecordsToDTO(value.ContextBefore), ContextAfter: logRecordsToDTO(value.ContextAfter),
		RawContextRedacted: value.RawContextRedacted,
	}
}

func logRecordsToDTO(values []domain.SanitizedLogRecord) []incidentdto.SanitizedLogRecordDTO {
	result := make([]incidentdto.SanitizedLogRecordDTO, 0, len(values))
	for _, value := range values {
		result = append(result, incidentdto.SanitizedLogRecordDTO{Timestamp: value.Timestamp, Stream: value.Stream, Message: value.Message, Redacted: value.Redacted})
	}
	return result
}
