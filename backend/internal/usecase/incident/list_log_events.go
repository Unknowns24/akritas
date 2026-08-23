package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (uc *UseCase) ListLogEvents(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.LogEvent], error) {
	if incidentID == uuid.Nil {
		return paging.Slice[domain.LogEvent]{}, domain.ErrIncidentNotFound
	}
	if _, err := uc.incidents.Get(ctx, incidentID); err != nil {
		return paging.Slice[domain.LogEvent]{}, err
	}
	return uc.incidents.ListLogEvents(ctx, incidentID, params)
}
