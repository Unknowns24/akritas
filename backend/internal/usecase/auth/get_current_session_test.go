package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestGetCurrentSessionHappyPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	admin, err := domain.NewAdministrator(uuid.New(), "admin@example.com", "Akritas Administrator", now.Add(-time.Hour))
	if err != nil {
		t.Fatalf("build administrator: %v", err)
	}
	session, err := domain.NewAdministratorSession(uuid.New(), admin.ID, now, now.Add(12*time.Hour), now.Add(168*time.Hour))
	if err != nil {
		t.Fatalf("build session: %v", err)
	}

	administrators := &fakeAdministratorRepository{findByIDResult: admin}
	uc := NewGetCurrentSessionUseCase(administrators)

	output, err := uc.Execute(context.Background(), *session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output.Administrator.ID != admin.ID {
		t.Fatal("output must expose the administrator found by session.AdministratorID")
	}
	if !output.AuthenticatedAt.Equal(session.AuthenticatedAt) ||
		!output.IdleExpiresAt.Equal(session.IdleExpiresAt) ||
		!output.AbsoluteExpiresAt.Equal(session.AbsoluteExpiresAt) {
		t.Fatal("output must reflect the timestamps of the session it was given, not recompute them")
	}
}

func TestGetCurrentSessionAdministratorNotFound(t *testing.T) {
	t.Parallel()

	session := domain.AdministratorSession{AdministratorID: uuid.New()}
	administrators := &fakeAdministratorRepository{findByIDResult: nil}
	uc := NewGetCurrentSessionUseCase(administrators)

	if _, err := uc.Execute(context.Background(), session); !errors.Is(err, domain.ErrInactiveAdministratorSession) {
		t.Fatalf("expected ErrInactiveAdministratorSession, got %v", err)
	}
}

func TestGetCurrentSessionFindByIDErrorPropagates(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("database unavailable")
	session := domain.AdministratorSession{AdministratorID: uuid.New()}
	administrators := &fakeAdministratorRepository{findByIDErr: wantErr}
	uc := NewGetCurrentSessionUseCase(administrators)

	if _, err := uc.Execute(context.Background(), session); !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
