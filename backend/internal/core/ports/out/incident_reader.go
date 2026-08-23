package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// IncidentReader is the minimal read boundary Investigation needs from H2
// (Detection + Incidents). H2 exists today only as a completed domain entity
// (internal/core/domain/incident.go); a real, PostgreSQL-backed
// implementation of this port is pending H2's merge. Get returns
// domain.ErrIncidentNotFound when the incident does not exist.
type IncidentReader interface {
	Exists(ctx context.Context, incidentID uuid.UUID) (bool, error)
	Get(ctx context.Context, incidentID uuid.UUID) (*domain.Incident, error)
}
