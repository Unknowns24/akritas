package incident

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type SanitizedLogRecordDTO struct {
	Timestamp time.Time        `json:"timestamp"`
	Stream    domain.LogStream `json:"stream,omitempty"`
	Message   string           `json:"message"`
	Redacted  bool             `json:"redacted"`
}
