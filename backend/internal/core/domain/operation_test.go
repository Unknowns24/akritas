package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewOperationStartsQueued(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	resourceType := OperationResourceInvestigation
	resourceID := uuid.New()
	key := "idem-key"
	operation, err := NewOperation(uuid.New(), OperationTypeInvestigation, &resourceType, &resourceID, &key, "queued for execution", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if operation.Status != OperationStatusQueued || operation.FinishedAt != nil {
		t.Fatal("new operation must start queued without a finish time")
	}
}

func TestOperationRejectsMismatchedResourcePair(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	resourceType := OperationResourceInvestigation
	if _, err := NewOperation(uuid.New(), OperationTypeInvestigation, &resourceType, nil, nil, "", now); !errors.Is(err, ErrInvalidOperation) {
		t.Fatalf("expected ErrInvalidOperation for mismatched resource pair, got %v", err)
	}
}

func TestOperationLifecycleSucceeds(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	operation, err := NewOperation(uuid.New(), OperationTypeInvestigation, nil, nil, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := operation.Start(now.Add(time.Second)); err != nil {
		t.Fatalf("start: %v", err)
	}
	if operation.Status != OperationStatusRunning {
		t.Fatal("expected running after Start")
	}
	if err := operation.Succeed(now.Add(2*time.Second), "done"); err != nil {
		t.Fatalf("succeed: %v", err)
	}
	if operation.Status != OperationStatusSucceeded || operation.FinishedAt == nil {
		t.Fatal("expected succeeded with a finish time")
	}
	if err := operation.Validate(); err != nil {
		t.Fatalf("terminal operation must validate: %v", err)
	}
}

func TestOperationLifecycleFails(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	operation, err := NewOperation(uuid.New(), OperationTypeInvestigation, nil, nil, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := operation.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	code := "0x504002N"
	if err := operation.Fail(now.Add(2*time.Second), "not implemented yet", &code); err != nil {
		t.Fatalf("fail: %v", err)
	}
	if operation.Status != OperationStatusFailed || operation.FailureCode == nil || *operation.FailureCode != code {
		t.Fatal("expected failed operation to carry the failure code")
	}
}

func TestOperationRejectsInvalidTransitions(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	operation, err := NewOperation(uuid.New(), OperationTypeInvestigation, nil, nil, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := operation.Succeed(now, "too early"); !errors.Is(err, ErrOperationTransition) {
		t.Fatalf("succeed before start must fail with ErrOperationTransition, got %v", err)
	}
	if err := operation.Fail(now, "too early", nil); !errors.Is(err, ErrOperationTransition) {
		t.Fatalf("fail before start must fail with ErrOperationTransition, got %v", err)
	}
	if err := operation.Start(now.Add(-time.Second)); !errors.Is(err, ErrOperationTransition) {
		t.Fatalf("start before creation time must fail, got %v", err)
	}
}
