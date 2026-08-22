package out

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type ProviderConnectionResult struct {
	Status      domain.ConnectionTestStatus
	CheckedAt   time.Time
	Latency     time.Duration
	UserMessage string
}
