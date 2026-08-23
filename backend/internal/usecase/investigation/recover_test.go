package investigation

import (
	"context"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestRecoverRequeuesOnlyPendingQueuedPairs(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	deps := newStartDeps()
	deps.now = now
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), now)
	if err != nil {
		t.Fatal(err)
	}
	resource := domain.OperationResourceInvestigation
	operation, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, &resource, &investigation.ID, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	deps.investigations.open = []domain.Investigation{*investigation}
	deps.operations.findByIDResult = operation
	if err := deps.runUseCase().Recover(context.Background(), deps.dispatcher); err != nil {
		t.Fatal(err)
	}
	if !deps.dispatcher.dispatched || deps.dispatcher.investigationID != investigation.ID {
		t.Fatal("pending/queued work was not requeued")
	}
}

func TestRecoverFailsInterruptedRunningWorkAndIncident(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	deps := newStartDeps()
	deps.now = now.Add(time.Minute)
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := investigation.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	resource := domain.OperationResourceInvestigation
	operation, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, &resource, &investigation.ID, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := operation.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	deps.incidents.getResult = &domain.Incident{ID: investigation.IncidentID, ProjectID: uuid.New(), Phase: domain.IncidentPhaseInvestigating}
	deps.investigations.open = []domain.Investigation{*investigation}
	deps.operations.findByIDResult = operation
	if err := deps.runUseCase().Recover(context.Background(), deps.dispatcher); err != nil {
		t.Fatal(err)
	}
	if deps.investigations.updated[len(deps.investigations.updated)-1].Status != domain.InvestigationStatusFailed ||
		deps.operations.updated[len(deps.operations.updated)-1].Status != domain.OperationStatusFailed ||
		deps.incidents.getResult.Phase != domain.IncidentPhaseFailed {
		t.Fatalf("running recovery did not fail durable state: investigation=%+v operation=%+v incident=%+v", deps.investigations.updated, deps.operations.updated, deps.incidents.getResult)
	}
}
