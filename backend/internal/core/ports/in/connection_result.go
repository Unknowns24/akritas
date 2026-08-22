package in

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

type ConnectionTestResult struct {
	Status      domain.ConnectionTestStatus
	CheckedAt   string
	LatencyMS   *int64
	UserMessage string
}
