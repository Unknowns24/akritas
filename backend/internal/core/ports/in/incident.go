package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type IncidentUseCase interface {
	Get(context.Context, uuid.UUID) (*domain.Incident, error)
	List(context.Context, paging.Params) (paging.Slice[domain.Incident], error)
	ListLogEvents(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.LogEvent], error)
}
