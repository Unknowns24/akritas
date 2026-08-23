// Package evidenceassembly implements out.EvidenceAssembler by combining
// existing read ports (IncidentReader, ProjectStore). It never invents
// evidence: a missing incident or project yields no evidence, not a
// placeholder.
package evidenceassembly

import (
	"context"
	"errors"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type Assembler struct {
	incidents portsout.IncidentReader
	projects  portsout.ProjectStore
	newID     func() uuid.UUID
	now       func() time.Time
}

func New(incidents portsout.IncidentReader, projects portsout.ProjectStore, newID func() uuid.UUID, now func() time.Time) portsout.EvidenceAssembler {
	return &Assembler{incidents: incidents, projects: projects, newID: newID, now: now}
}

func (a *Assembler) Assemble(ctx context.Context, investigation domain.Investigation) ([]domain.Evidence, error) {
	incident, err := a.incidents.Get(ctx, investigation.IncidentID)
	if errors.Is(err, domain.ErrIncidentNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	project, err := a.projects.Get(ctx, incident.ProjectID)
	if errors.Is(err, domain.ErrProjectNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	deploymentMetadata, err := deploymentMetadataEvidence(a.newID(), investigation.ID, *project, a.now().UTC())
	if err != nil {
		return nil, err
	}
	return []domain.Evidence{*deploymentMetadata}, nil
}
