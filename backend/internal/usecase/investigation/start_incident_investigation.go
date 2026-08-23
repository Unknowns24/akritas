package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
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

	var operation *domain.Operation
	var investigationID uuid.UUID
	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		replayed, findErr := uc.operations.FindByIdempotencyKey(txCtx, idempotencyKey)
		if findErr != nil {
			return findErr
		}
		if replayed != nil {
			operation = replayed
			return nil
		}

		incident, lockErr := uc.incidents.Lock(txCtx, command.IncidentID)
		if lockErr != nil {
			return lockErr
		}
		active, activeErr := uc.investigations.ExistsActiveForIncident(txCtx, command.IncidentID)
		if activeErr != nil {
			return activeErr
		}
		if active {
			return domain.ErrInvestigationAlreadyActive
		}
		if transitionErr := incident.StartInvestigation(); transitionErr != nil {
			return transitionErr
		}

		now := uc.now().UTC()
		created, createErr := domain.NewInvestigation(uc.newID(), command.IncidentID, now)
		if createErr != nil {
			return createErr
		}
		if createErr = uc.investigations.Create(txCtx, created); createErr != nil {
			return createErr
		}

		resourceType := domain.OperationResourceInvestigation
		resourceID := created.ID
		operation, createErr = domain.NewOperation(
			uc.newID(), domain.OperationTypeInvestigation, &resourceType, &resourceID,
			&idempotencyKey, "La investigación fue encolada.", now,
		)
		if createErr != nil {
			return createErr
		}
		if createErr = uc.operations.Create(txCtx, operation); createErr != nil {
			return createErr
		}
		if updateErr := uc.incidents.Update(txCtx, incident); updateErr != nil {
			return updateErr
		}
		investigationID = created.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	if investigationID != uuid.Nil {
		uc.dispatcher.Dispatch(investigationID, operation.ID)
	}
	return operation, nil
}
