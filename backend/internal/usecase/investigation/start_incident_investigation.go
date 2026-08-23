package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func (uc *UseCase) StartIncidentInvestigation(ctx context.Context, command portsin.StartIncidentInvestigationCommand) (*domain.Operation, error) {
	idempotencyKey := command.IdempotencyKey.String()
	existing, err := uc.operations.FindByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	exists, err := uc.incidents.Exists(ctx, command.IncidentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrIncidentNotFound
	}

	active, err := uc.investigations.ExistsActiveForIncident(ctx, command.IncidentID)
	if err != nil {
		return nil, err
	}
	if active {
		return nil, domain.ErrInvestigationAlreadyActive
	}

	now := uc.now().UTC()
	created, err := domain.NewInvestigation(uc.newID(), command.IncidentID, now)
	if err != nil {
		return nil, err
	}
	if err := uc.investigations.Create(ctx, created); err != nil {
		return nil, err
	}

	resourceType := domain.OperationResourceInvestigation
	resourceID := created.ID
	operation, err := domain.NewOperation(
		uc.newID(), domain.OperationTypeInvestigation, &resourceType, &resourceID,
		&idempotencyKey, "La investigación fue encolada.", now,
	)
	if err != nil {
		return nil, err
	}
	if err := uc.operations.Create(ctx, operation); err != nil {
		return nil, err
	}

	uc.dispatcher.Dispatch(created.ID, operation.ID)
	return operation, nil
}
