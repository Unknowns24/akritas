package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestLogoutAdministratorHappyPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session, err := domain.NewAdministratorSession(uuid.New(), uuid.New(), now.Add(-time.Hour), now.Add(11*time.Hour), now.Add(168*time.Hour))
	if err != nil {
		t.Fatalf("build session: %v", err)
	}

	sessions := &fakeAdministratorSessionRepository{}
	clock := &fakeClock{now: now}
	uc := NewLogoutAdministratorUseCase(sessions, clock)

	if err := uc.Execute(context.Background(), *session); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sessions.revokeCalled || sessions.revokedSessionID != session.ID {
		t.Fatal("Revoke must be called for the session's id")
	}
	if !sessions.revokedAt.Equal(now) {
		t.Fatalf("revokedAt = %v, want %v", sessions.revokedAt, now)
	}
}

func TestLogoutAdministratorRepositoryErrorPropagates(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session, err := domain.NewAdministratorSession(uuid.New(), uuid.New(), now.Add(-time.Hour), now.Add(11*time.Hour), now.Add(168*time.Hour))
	if err != nil {
		t.Fatalf("build session: %v", err)
	}

	wantErr := errors.New("database unavailable")
	sessions := &fakeAdministratorSessionRepository{revokeErr: wantErr}
	uc := NewLogoutAdministratorUseCase(sessions, &fakeClock{now: now})

	if err := uc.Execute(context.Background(), *session); !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
