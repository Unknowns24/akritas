package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type InvestigationGetter interface {
	FindByID(context.Context, uuid.UUID) (*domain.Investigation, error)
}

type LatestInvestigationFinder interface {
	FindLatestByIncident(context.Context, uuid.UUID) (*domain.Investigation, error)
}

type InvestigationStore interface {
	InvestigationGetter
	LatestInvestigationFinder
	Create(context.Context, *domain.Investigation) error
	Update(context.Context, *domain.Investigation) error
	ListByIncident(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.Investigation], error)
	ExistsActiveForIncident(context.Context, uuid.UUID) (bool, error)
	ListOpen(context.Context) ([]domain.Investigation, error)
}
