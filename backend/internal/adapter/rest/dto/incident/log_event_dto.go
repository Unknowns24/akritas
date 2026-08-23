package incident

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type LogEventDTO struct {
	ID                 uuid.UUID               `json:"id"`
	ProjectID          uuid.UUID               `json:"project_id"`
	Timestamp          time.Time               `json:"timestamp"`
	Severity           domain.Severity         `json:"severity"`
	Message            string                  `json:"message"`
	Fingerprint        string                  `json:"fingerprint"`
	DetectionRules     []string                `json:"detection_rules"`
	ContextBefore      []SanitizedLogRecordDTO `json:"context_before,omitempty"`
	ContextAfter       []SanitizedLogRecordDTO `json:"context_after,omitempty"`
	RawContextRedacted bool                    `json:"raw_context_redacted"`
}
