package out

import (
	"context"

	"github.com/google/uuid"
)

// IncidentReader is the minimal read boundary Investigation needs from H2
// (Detection + Incidents) to validate an incident_id. H2 exists today only as
// a completed domain entity (internal/core/domain/incident.go); a real,
// PostgreSQL-backed implementation of this port is pending H2's merge.
type IncidentReader interface {
	Exists(ctx context.Context, incidentID uuid.UUID) (bool, error)
}
