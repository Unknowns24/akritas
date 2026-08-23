package operation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestRepositoryPersistsCreatesFindsAndUpdates(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)

	resourceType := domain.OperationResourceInvestigation
	resourceID := uuid.New()
	key := uuid.New().String()
	value, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, &resourceType, &resourceID, &key, "encolada", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, value); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repository.FindByID(ctx, value.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found.Status != domain.OperationStatusQueued || found.ResourceID == nil || *found.ResourceID != resourceID {
		t.Fatalf("expected a queued operation with the resource pair persisted, got %+v", found)
	}

	if err := found.Start(now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := repository.Update(ctx, found); err != nil {
		t.Fatalf("update to running: %v", err)
	}
	if err := found.Succeed(now.Add(2*time.Second), "listo"); err != nil {
		t.Fatal(err)
	}
	if err := repository.Update(ctx, found); err != nil {
		t.Fatalf("update to succeeded: %v", err)
	}

	reloaded, err := repository.FindByID(ctx, value.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Status != domain.OperationStatusSucceeded || reloaded.FinishedAt == nil || reloaded.IdempotencyKey == nil || *reloaded.IdempotencyKey != key {
		t.Fatalf("expected the persisted operation to be succeeded and keep its idempotency key, got %+v", reloaded)
	}
}

func TestRepositoryFindByIDNotFound(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.FindByID(context.Background(), uuid.New()); !errors.Is(err, domain.ErrOperationNotFound) {
		t.Fatalf("expected ErrOperationNotFound, got %v", err)
	}
}

func TestRepositoryFindByIdempotencyKeyHitAndMiss(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	key := uuid.New().String()

	value, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, nil, nil, &key, "encolada", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, value); err != nil {
		t.Fatal(err)
	}

	found, err := repository.FindByIdempotencyKey(ctx, key)
	if err != nil {
		t.Fatalf("find by key: %v", err)
	}
	if found == nil || found.ID != value.ID {
		t.Fatal("expected the operation tied to the idempotency key to be found")
	}

	miss, err := repository.FindByIdempotencyKey(ctx, uuid.New().String())
	if err != nil || miss != nil {
		t.Fatalf("expected a miss to return (nil, nil), got %+v, %v", miss, err)
	}
}

func TestRepositoryAllowsMultipleOperationsWithoutIdempotencyKey(t *testing.T) {
	db := dbtest.Connect(t)
	repository, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)

	first, err := domain.NewOperation(uuid.New(), domain.OperationTypeSystemDiagnostics, nil, nil, nil, "", now)
	if err != nil {
		t.Fatal(err)
	}
	second, err := domain.NewOperation(uuid.New(), domain.OperationTypeSystemDiagnostics, nil, nil, nil, "", now)
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.Create(ctx, first); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if err := repository.Create(ctx, second); err != nil {
		t.Fatalf("create second without idempotency key must not collide: %v", err)
	}
}
