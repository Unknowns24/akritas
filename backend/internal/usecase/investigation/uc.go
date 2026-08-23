package investigation

import (
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

// UseCase satisfies both portsin.InvestigationUseCase (called by REST) and
// portsin.RunInvestigationUseCase (called only by InvestigationDispatcher).
type UseCase struct {
	incidents      portsout.IncidentReader
	investigations portsout.InvestigationStore
	operations     portsout.OperationStore
	dispatcher     portsout.InvestigationDispatcher
	runner         portsout.InvestigationRunner
	newID          func() uuid.UUID
	now            func() time.Time
}

func New(
	incidents portsout.IncidentReader,
	investigations portsout.InvestigationStore,
	operations portsout.OperationStore,
	dispatcher portsout.InvestigationDispatcher,
	runner portsout.InvestigationRunner,
	newID func() uuid.UUID,
	now func() time.Time,
) *UseCase {
	return &UseCase{
		incidents: incidents, investigations: investigations, operations: operations,
		dispatcher: dispatcher, runner: runner, newID: newID, now: now,
	}
}
