package incident

import "time"

type TimelineEventDTO struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	Type       string    `json:"type"`
	OccurredAt time.Time `json:"occurred_at"`
	Summary    string    `json:"summary"`
	Detail     string    `json:"detail,omitempty"`
}
