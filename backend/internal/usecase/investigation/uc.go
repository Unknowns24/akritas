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
	incidents      portsout.IncidentReader
	investigations portsout.InvestigationStore
	operations     portsout.OperationStore
	dispatcher     portsout.InvestigationDispatcher
	newID          func() uuid.UUID
	now            func() time.Time
}

func New(
	incidents portsout.IncidentReader,
	investigations portsout.InvestigationStore,
	operations portsout.OperationStore,
	dispatcher portsout.InvestigationDispatcher,
	newID func() uuid.UUID,
	now func() time.Time,
) *UseCase {
	return &UseCase{
		incidents: incidents, investigations: investigations, operations: operations,
		dispatcher: dispatcher, newID: newID, now: now,
	}
}

// RunUseCase implements portsin.RunInvestigationUseCase (Execute, called only
// by InvestigationDispatcher). It never depends on InvestigationDispatcher:
// the dispatcher is built from this value, so depending on the dispatcher
// back would be a construction cycle.
type RunUseCase struct {
	investigations portsout.InvestigationStore
	operations     portsout.OperationStore
	evidence       portsout.EvidenceStore
	assembler      portsout.EvidenceAssembler
	runner         portsout.InvestigationRunner
	now            func() time.Time
}

func NewRunUseCase(
	investigations portsout.InvestigationStore,
	operations portsout.OperationStore,
	evidence portsout.EvidenceStore,
	assembler portsout.EvidenceAssembler,
	runner portsout.InvestigationRunner,
	now func() time.Time,
) *RunUseCase {
	return &RunUseCase{
		investigations: investigations, operations: operations,
		evidence: evidence, assembler: assembler, runner: runner, now: now,
	}
}
