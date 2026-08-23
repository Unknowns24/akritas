package mapper

import (
	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func ConnectionTestToDTO(value portsin.ConnectionTestResult) commondto.ConnectionTestDTO {
	return commondto.ConnectionTestDTO{Status: string(value.Status), CheckedAt: value.CheckedAt, LatencyMS: value.LatencyMS, UserMessage: value.UserMessage}
}
