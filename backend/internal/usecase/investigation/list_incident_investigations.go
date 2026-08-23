package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (uc *UseCase) ListIncidentInvestigations(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.Investigation], error) {
	exists, err := uc.incidents.Exists(ctx, incidentID)
	if err != nil {
		return paging.Slice[domain.Investigation]{}, err
	}
	if !exists {
		return paging.Slice[domain.Investigation]{}, domain.ErrIncidentNotFound
	}
	return uc.investigations.ListByIncident(ctx, incidentID, params)
}
