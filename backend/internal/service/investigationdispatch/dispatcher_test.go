package investigationdispatch

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeRunUseCase struct {
	calls chan [2]uuid.UUID
}

func (f *fakeRunUseCase) Execute(ctx context.Context, investigationID, operationID uuid.UUID) error {
	f.calls <- [2]uuid.UUID{investigationID, operationID}
	return nil
}

func TestDispatchInvokesRunAsynchronously(t *testing.T) {
	t.Parallel()

	fake := &fakeRunUseCase{calls: make(chan [2]uuid.UUID, 1)}
	dispatcher := New(fake, time.Second)

	investigationID, operationID := uuid.New(), uuid.New()
	dispatcher.Dispatch(investigationID, operationID)

	select {
	case call := <-fake.calls:
		if call[0] != investigationID || call[1] != operationID {
			t.Fatal("expected Execute to receive the dispatched investigation and operation IDs")
		}
	case <-time.After(time.Second):
		t.Fatal("expected Execute to be invoked asynchronously")
	}
}
