package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

type DashboardUseCase interface {
	GetOverview(context.Context) (domain.Overview, error)
	ListActivity(context.Context, paging.Params) (paging.Slice[domain.ActivityEvent], error)
}
