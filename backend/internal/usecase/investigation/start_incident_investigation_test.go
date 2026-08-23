package investigation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

type startDeps struct {
	incidents      *fakeIncidentReader
	investigations *fakeInvestigationStore
	operations     *fakeOperationStore
	dispatcher     *fakeInvestigationDispatcher
	runner         *fakeInvestigationRunner
	now            time.Time
}

func newStartDeps() *startDeps {
	return &startDeps{
		incidents:      &fakeIncidentReader{exists: true},
		investigations: &fakeInvestigationStore{},
		operations:     &fakeOperationStore{},
		dispatcher:     &fakeInvestigationDispatcher{},
		runner:         &fakeInvestigationRunner{},
		now:            time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
	}
}

func (d *startDeps) usecase() *UseCase {
	id := uuid.New()
	ids := []uuid.UUID{id, uuid.New()}
	next := 0
	newID := func() uuid.UUID {
		if next >= len(ids) {
			return uuid.New()
		}
		value := ids[next]
		next++
		return value
	}
	return New(d.incidents, d.investigations, d.operations, d.dispatcher, newID, func() time.Time { return d.now })
}

func (d *startDeps) runUseCase() *RunUseCase {
	return NewRunUseCase(d.investigations, d.operations, d.runner, func() time.Time { return d.now })
}

func TestStartIncidentInvestigationHappyPathQueuesAndDispatches(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	incidentID := uuid.New()
	operation, err := deps.usecase().StartIncidentInvestigation(context.Background(), portsin.StartIncidentInvestigationCommand{IncidentID: incidentID, IdempotencyKey: uuid.New()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deps.investigations.created == nil || deps.investigations.created.Status != domain.InvestigationStatusPending {
		t.Fatal("expected a pending investigation to be created")
	}
	if deps.operations.created == nil || deps.operations.created.Status != domain.OperationStatusQueued {
		t.Fatal("expected a queued operation to be created")
	}
	if operation == nil || operation.ID != deps.operations.created.ID {
		t.Fatal("expected the queued operation to be returned")
	}
	if !deps.dispatcher.dispatched || deps.dispatcher.investigationID != deps.investigations.created.ID || deps.dispatcher.operationID != deps.operations.created.ID {
		t.Fatal("expected the dispatcher to be invoked with the new investigation and operation")
	}
}

func TestStartIncidentInvestigationRejectsMissingIncident(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	deps.incidents.exists = false
	_, err := deps.usecase().StartIncidentInvestigation(context.Background(), portsin.StartIncidentInvestigationCommand{IncidentID: uuid.New(), IdempotencyKey: uuid.New()})
	if !errors.Is(err, domain.ErrIncidentNotFound) {
		t.Fatalf("expected ErrIncidentNotFound, got %v", err)
	}
	if deps.investigations.created != nil || deps.operations.created != nil || deps.dispatcher.dispatched {
		t.Fatal("nothing must be created or dispatched when the incident does not exist")
	}
}

func TestStartIncidentInvestigationRejectsDuplicateActive(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	deps.investigations.activeForIncident = true
	_, err := deps.usecase().StartIncidentInvestigation(context.Background(), portsin.StartIncidentInvestigationCommand{IncidentID: uuid.New(), IdempotencyKey: uuid.New()})
	if !errors.Is(err, domain.ErrInvestigationAlreadyActive) {
		t.Fatalf("expected ErrInvestigationAlreadyActive, got %v", err)
	}
	if deps.investigations.created != nil || deps.operations.created != nil || deps.dispatcher.dispatched {
		t.Fatal("nothing must be created or dispatched when an investigation is already active")
	}
}

func TestStartIncidentInvestigationReplaysIdempotencyKeyWithoutCreating(t *testing.T) {
	t.Parallel()
	deps := newStartDeps()
	existing, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, nil, nil, nil, "ya encolada", deps.now)
	if err != nil {
		t.Fatal(err)
	}
	deps.operations.findByKeyResult = existing
	operation, err := deps.usecase().StartIncidentInvestigation(context.Background(), portsin.StartIncidentInvestigationCommand{IncidentID: uuid.New(), IdempotencyKey: uuid.New()})
	if err != nil {
		t.Fatal(err)
	}
	if operation.ID != existing.ID {
		t.Fatal("replay must return the operation already tied to the idempotency key")
	}
	if !deps.operations.findByKeyCalled {
		t.Fatal("expected the idempotency key to be looked up")
	}
	if deps.investigations.created != nil || deps.operations.created != nil || deps.dispatcher.dispatched {
		t.Fatal("a replay must not create a new investigation, operation or dispatch")
	}
}
