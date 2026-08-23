package detection

import (
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

const maximumLogicalEventBytes = 20000

type LogicalEvent struct {
	Records []domain.SanitizedLogRecord
}

func (e LogicalEvent) Message() string {
	parts := make([]string, 0, len(e.Records))
	length := 0
	for _, record := range e.Records {
		message := strings.TrimSpace(record.Message)
		if message == "" {
			continue
		}
		remaining := maximumLogicalEventBytes - length
		if len(parts) > 0 {
			remaining--
		}
		if remaining <= 0 {
			break
		}
		if len(message) > remaining {
			message = message[:remaining]
		}
		parts = append(parts, message)
		length += len(message)
		if len(parts) > 1 {
			length++
		}
	}
	return strings.Join(parts, "\n")
}
