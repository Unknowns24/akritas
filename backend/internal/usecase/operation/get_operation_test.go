package operation

import (
	"context"
	"errors"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type fakeOperationStore struct {
	result *domain.Operation
	err    error
}

func (f *fakeOperationStore) Create(ctx context.Context, value *domain.Operation) error { return nil }
func (f *fakeOperationStore) Update(ctx context.Context, value *domain.Operation) error { return nil }
func (f *fakeOperationStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Operation, error) {
	return f.result, f.err
}
func (f *fakeOperationStore) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Operation, error) {
	return nil, nil
}

func TestGetOperationHappyPath(t *testing.T) {
	t.Parallel()
	store := &fakeOperationStore{result: &domain.Operation{ID: uuid.New()}}
	uc := New(store)
	value, err := uc.GetOperation(context.Background(), store.result.ID)
	if err != nil {
		t.Fatal(err)
	}
	if value.ID != store.result.ID {
		t.Fatal("expected the store's operation to be returned unchanged")
	}
}

func TestGetOperationNotFound(t *testing.T) {
	t.Parallel()
	store := &fakeOperationStore{err: domain.ErrOperationNotFound}
	uc := New(store)
	if _, err := uc.GetOperation(context.Background(), uuid.New()); !errors.Is(err, domain.ErrOperationNotFound) {
		t.Fatalf("expected ErrOperationNotFound, got %v", err)
	}
}
