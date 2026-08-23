package incident

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type IncidentDTO struct {
	ID               uuid.UUID                `json:"id"`
	Key              string                   `json:"key"`
	Project          ProjectReferenceDTO      `json:"project"`
	Fingerprint      string                   `json:"fingerprint"`
	Severity         domain.Severity          `json:"severity"`
	Title            string                   `json:"title"`
	Summary          string                   `json:"summary,omitempty"`
	Phase            domain.IncidentPhase     `json:"phase"`
	TerminalOutcome  *domain.TerminalOutcome  `json:"terminal_outcome,omitempty"`
	RootCauseStatus  *domain.RootCauseStatus  `json:"root_cause_status,omitempty"`
	ResolutionStatus *domain.ResolutionStatus `json:"resolution_status,omitempty"`
	Confidence       *float64                 `json:"confidence,omitempty"`
	OccurrenceCount  int64                    `json:"occurrence_count"`
	FirstSeenAt      time.Time                `json:"first_seen_at"`
	LastSeenAt       time.Time                `json:"last_seen_at"`
}
