package investigation

import (
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

// UseCase implements portsin.InvestigationUseCase (Start/Get/List, called by
// REST). It depends on InvestigationDispatcher but never on
// InvestigationRunner: that split avoids a construction cycle with RunUseCase
// below, which the dispatcher wraps.
type UseCase struct {
	incidents      portsout.IncidentWorkflowStore
	investigations portsout.InvestigationStore
	operations     portsout.OperationStore
	transactor     portsout.Transactor
	dispatcher     portsout.InvestigationDispatcher
	newID          func() uuid.UUID
	now            func() time.Time
}

func New(
	incidents portsout.IncidentWorkflowStore,
	investigations portsout.InvestigationStore,
	operations portsout.OperationStore,
	transactor portsout.Transactor,
	dispatcher portsout.InvestigationDispatcher,
	newID func() uuid.UUID,
	now func() time.Time,
) *UseCase {
	return &UseCase{
		incidents: incidents, investigations: investigations, operations: operations, transactor: transactor,
		dispatcher: dispatcher, newID: newID, now: now,
	}
}

// RunUseCase implements portsin.RunInvestigationUseCase (Execute, called only
// by InvestigationDispatcher). It never depends on InvestigationDispatcher:
// the dispatcher is built from this value, so depending on the dispatcher
// back would be a construction cycle.
type RunUseCase struct {
	incidents      portsout.IncidentWorkflowStore
	investigations portsout.InvestigationStore
	operations     portsout.OperationStore
	evidence       portsout.EvidenceStore
	assembler      portsout.InvestigationContextAssembler
	runner         portsout.InvestigationRunner
	transactor     portsout.Transactor
	now            func() time.Time
}

func NewRunUseCase(
	incidents portsout.IncidentWorkflowStore,
	investigations portsout.InvestigationStore,
	operations portsout.OperationStore,
	evidence portsout.EvidenceStore,
	assembler portsout.InvestigationContextAssembler,
	runner portsout.InvestigationRunner,
	transactor portsout.Transactor,
	now func() time.Time,
) *RunUseCase {
	return &RunUseCase{
		incidents: incidents, investigations: investigations, operations: operations,
		evidence: evidence, assembler: assembler, runner: runner, transactor: transactor, now: now,
	}
}
