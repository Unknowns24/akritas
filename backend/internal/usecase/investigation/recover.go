package investigation

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

const interruptedInvestigationMessage = "La investigación fue interrumpida durante el reinicio y puede reintentarse."

// Recover reconciles durable work for the documented single-instance MVP.
// Pending/queued pairs are safe to requeue; anything that had started is
// failed because an external QVAC/tool side effect cannot be resumed exactly.
func (uc *RunUseCase) Recover(ctx context.Context, dispatcher portsout.InvestigationDispatcher) error {
	open, err := uc.investigations.ListOpen(ctx)
	if err != nil {
		return err
	}
	for index := range open {
		investigation := open[index]
		operation, operationErr := uc.operations.FindByResource(ctx, domain.OperationResourceInvestigation, investigation.ID)
		if investigation.Status == domain.InvestigationStatusPending && operationErr == nil && operation.Status == domain.OperationStatusQueued {
			dispatcher.Dispatch(investigation.ID, operation.ID)
			continue
		}
		if operationErr != nil && !errors.Is(operationErr, domain.ErrOperationNotFound) {
			return operationErr
		}
		if err := uc.failInterrupted(ctx, &investigation, operation); err != nil {
			return err
		}
	}
	return nil
}

func (uc *RunUseCase) failInterrupted(ctx context.Context, investigation *domain.Investigation, operation *domain.Operation) error {
	at := uc.now().UTC()
	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		incident, err := uc.incidents.Lock(txCtx, investigation.IncidentID)
		if err != nil {
			return err
		}
		if investigation.Status == domain.InvestigationStatusPending {
			if err = investigation.Start(at); err != nil {
				return err
			}
		}
		if err = investigation.Fail(at, interruptedInvestigationMessage); err != nil {
			return err
		}
		if operation != nil {
			if operation.Status == domain.OperationStatusQueued {
				if err = operation.Start(at); err != nil {
					return err
				}
			}
			if operation.Status == domain.OperationStatusRunning {
				if err = operation.Fail(at, interruptedInvestigationMessage, nil); err != nil {
					return err
				}
			}
			if err = uc.operations.Update(txCtx, operation); err != nil {
				return err
			}
		}
		if incident.Phase == domain.IncidentPhaseInvestigating {
			if err = incident.FailInvestigation(); err != nil {
				return err
			}
			if err = uc.incidents.Update(txCtx, incident); err != nil {
				return err
			}
		}
		return uc.investigations.Update(txCtx, investigation)
	})
}
