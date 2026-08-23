package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (uc *UseCase) ListTimeline(ctx context.Context, id uuid.UUID, params paging.Params) (paging.Slice[domain.TimelineEvent], error) {
	if id == uuid.Nil || uc.timeline == nil {
		return paging.Slice[domain.TimelineEvent]{}, domain.ErrIncidentNotFound
	}
	if _, err := uc.incidents.Get(ctx, id); err != nil {
		return paging.Slice[domain.TimelineEvent]{}, err
	}
	return uc.timeline.ListTimeline(ctx, id, params)
}
