package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (uc *UseCase) ListApplications(ctx context.Context, id uuid.UUID, params paging.Params) (paging.Slice[domain.DokployApplication], error) {
	server, err := uc.store.Get(ctx, id)
	if err != nil {
		return paging.Slice[domain.DokployApplication]{}, err
	}
	page, err := uc.gateway.ListApplications(ctx, *server, params)
	if err != nil {
		return paging.Slice[domain.DokployApplication]{}, err
	}
	now := uc.now().UTC()
	server.ApplicationCount = int(page.Total)
	server.LastSyncedAt = &now
	server.UpdatedAt = now
	if err := uc.store.UpdateConnection(ctx, server); err != nil {
		return paging.Slice[domain.DokployApplication]{}, err
	}
	return page, nil
}
