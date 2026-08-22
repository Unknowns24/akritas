package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func ConnectionTestToDTO(value portsin.ConnectionTestResult) dto.ConnectionTestDTO {
	return dto.ConnectionTestDTO{Status: string(value.Status), CheckedAt: value.CheckedAt, LatencyMS: value.LatencyMS, UserMessage: value.UserMessage}
}
