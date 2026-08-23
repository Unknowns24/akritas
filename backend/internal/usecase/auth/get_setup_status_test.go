package auth

import (
	"context"
	"errors"
	"testing"
)

func TestGetSetupStatusExecute(t *testing.T) {
	t.Parallel()

	t.Run("no administrator yet", func(t *testing.T) {
		t.Parallel()
		administrators := &fakeAdministratorRepository{exists: false}
		uc := NewGetSetupStatusUseCase(administrators)

		status, err := uc.Execute(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !status.SetupRequired || !status.RegistrationOpen {
			t.Fatalf("expected setup_required and registration_open true, got %+v", status)
		}
	})

	t.Run("administrator already exists", func(t *testing.T) {
		t.Parallel()
		administrators := &fakeAdministratorRepository{exists: true}
		uc := NewGetSetupStatusUseCase(administrators)

		status, err := uc.Execute(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status.SetupRequired || status.RegistrationOpen {
			t.Fatalf("expected setup_required and registration_open false, got %+v", status)
		}
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()
		repoErr := errors.New("connection lost")
		administrators := &fakeAdministratorRepository{err: repoErr}
		uc := NewGetSetupStatusUseCase(administrators)

		if _, err := uc.Execute(context.Background()); !errors.Is(err, repoErr) {
			t.Fatalf("expected repository error, got %v", err)
		}
	})
}
