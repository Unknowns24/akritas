package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

func (uc *UseCase) List(ctx context.Context, params paging.Params) (paging.Slice[domain.DokployServer], error) {
	return uc.store.List(ctx, params)
}
