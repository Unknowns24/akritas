package dokployserver

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (uc *UseCase) ListComposes(ctx context.Context, id uuid.UUID, params paging.Params) (paging.Slice[domain.DokployCompose], error) {
	server, err := uc.store.Get(ctx, id)
	if err != nil {
		return paging.Slice[domain.DokployCompose]{}, err
	}
	page, err := uc.gateway.ListComposes(ctx, *server, params)
	if err != nil {
		return paging.Slice[domain.DokployCompose]{}, err
	}
	now := uc.now().UTC()
	server.ComposeCount = int(page.Total)
	server.LastSyncedAt = &now
	server.UpdatedAt = now
	if err := uc.store.UpdateConnection(ctx, server); err != nil {
		return paging.Slice[domain.DokployCompose]{}, err
	}
	return page, nil
}

func (uc *UseCase) ListComposeServices(ctx context.Context, serverID uuid.UUID, composeID string, refresh bool) ([]domain.DokployComposeService, error) {
	composeID = strings.TrimSpace(composeID)
	if serverID == uuid.Nil || composeID == "" {
		return nil, domain.ErrInvalidDokployCompose
	}
	server, err := uc.store.Get(ctx, serverID)
	if err != nil {
		return nil, err
	}
	return uc.gateway.ListComposeServices(ctx, *server, composeID, refresh)
}
