package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type StartIncidentInvestigationCommand struct {
	IncidentID     uuid.UUID
	IdempotencyKey uuid.UUID
}

type InvestigationUseCase interface {
	StartIncidentInvestigation(context.Context, StartIncidentInvestigationCommand) (*domain.Operation, error)
	GetInvestigation(context.Context, uuid.UUID) (*domain.Investigation, error)
	ListIncidentInvestigations(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.Investigation], error)
}

// RunInvestigationUseCase orchestrates the asynchronous pending->running->
// completed/failed transition. It is not part of InvestigationUseCase because
// REST never calls it directly: only InvestigationDispatcher does, after
// StartIncidentInvestigation has already returned its 202 response.
type RunInvestigationUseCase interface {
	Execute(ctx context.Context, investigationID, operationID uuid.UUID) error
}
