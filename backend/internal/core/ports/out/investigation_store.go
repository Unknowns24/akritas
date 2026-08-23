package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type InvestigationStore interface {
	Create(context.Context, *domain.Investigation) error
	Update(context.Context, *domain.Investigation) error
	FindByID(context.Context, uuid.UUID) (*domain.Investigation, error)
	ListByIncident(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.Investigation], error)
	ExistsActiveForIncident(context.Context, uuid.UUID) (bool, error)
}
