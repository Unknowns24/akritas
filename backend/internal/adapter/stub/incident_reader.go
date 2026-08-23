// Package stub holds safe-by-default placeholders for ports whose real
// implementation depends on a parallel workstream that has not merged yet.
// Nothing here should outlive the workstream it stands in for.
package stub

import (
	"context"

	"github.com/google/uuid"
)

// DenyAllIncidentReader stands in for H2's real IncidentReader until H2
// (Detection + Incidents) merges a PostgreSQL-backed implementation. It
// always reports that the incident does not exist, so investigation routes
// stay live and auditable in production without ever over-authorizing:
// creating an Investigation will 404 until H2 lands.
type DenyAllIncidentReader struct{}

func NewDenyAllIncidentReader() *DenyAllIncidentReader {
	return &DenyAllIncidentReader{}
}

func (r *DenyAllIncidentReader) Exists(ctx context.Context, incidentID uuid.UUID) (bool, error) {
	return false, nil
}
