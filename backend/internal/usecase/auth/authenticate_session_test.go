package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

const authSessionIdleTTL = 12 * time.Hour

func newTestSession(t *testing.T, now time.Time) *domain.AdministratorSession {
	t.Helper()
	session, err := domain.NewAdministratorSession(
		uuid.New(), uuid.New(), now.Add(-time.Hour), now.Add(11*time.Hour), now.Add(168*time.Hour),
	)
	if err != nil {
		t.Fatalf("build session: %v", err)
	}
	return session
}

func TestAuthenticateSessionHappyPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session := newTestSession(t, now)
	sessions := &fakeAdministratorSessionRepository{findByTokenHashResult: session}
	sessionTokens := &fakeSessionTokenGenerator{}
	clock := &fakeClock{now: now}
	uc := NewAuthenticateSessionUseCase(sessions, sessionTokens, clock.Now, authSessionIdleTTL)

	resolved, err := uc.Execute(context.Background(), "raw-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.ID != session.ID {
		t.Fatal("resolved session must be the one found by token hash")
	}
	if !resolved.IdleExpiresAt.Equal(now.Add(authSessionIdleTTL)) {
		t.Fatalf("IdleExpiresAt = %v, want %v (extended)", resolved.IdleExpiresAt, now.Add(authSessionIdleTTL))
	}
	if !sessions.updateIdleCalled || sessions.updatedIdleSessionID != session.ID {
		t.Fatal("UpdateIdleExpiry must be called for the resolved session")
	}
	if !sessions.updatedIdleExpiresAt.Equal(now.Add(authSessionIdleTTL)) {
		t.Fatalf("persisted idle expiry = %v, want %v", sessions.updatedIdleExpiresAt, now.Add(authSessionIdleTTL))
	}
}

func TestAuthenticateSessionEmptyToken(t *testing.T) {
	t.Parallel()

	sessions := &fakeAdministratorSessionRepository{}
	sessionTokens := &fakeSessionTokenGenerator{}
	clock := &fakeClock{now: time.Now()}
	uc := NewAuthenticateSessionUseCase(sessions, sessionTokens, clock.Now, authSessionIdleTTL)

	_, err := uc.Execute(context.Background(), "")
	if !errors.Is(err, domain.ErrInactiveAdministratorSession) {
		t.Fatalf("expected ErrInactiveAdministratorSession, got %v", err)
	}
	if sessions.findByTokenHashCalled {
		t.Fatal("no lookup should happen for an empty token")
	}
}

func TestAuthenticateSessionNotFound(t *testing.T) {
	t.Parallel()

	sessions := &fakeAdministratorSessionRepository{findByTokenHashResult: nil}
	sessionTokens := &fakeSessionTokenGenerator{}
	clock := &fakeClock{now: time.Now()}
	uc := NewAuthenticateSessionUseCase(sessions, sessionTokens, clock.Now, authSessionIdleTTL)

	_, err := uc.Execute(context.Background(), "raw-token")
	if !errors.Is(err, domain.ErrInactiveAdministratorSession) {
		t.Fatalf("expected ErrInactiveAdministratorSession, got %v", err)
	}
	if sessions.updateIdleCalled {
		t.Fatal("no extension should happen when the session is not found")
	}
}

func TestAuthenticateSessionExpiredOrRevoked(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	t.Run("expired", func(t *testing.T) {
		t.Parallel()
		session, err := domain.NewAdministratorSession(uuid.New(), uuid.New(), now.Add(-2*time.Hour), now.Add(-time.Hour), now.Add(168*time.Hour))
		if err != nil {
			t.Fatalf("build session: %v", err)
		}
		sessions := &fakeAdministratorSessionRepository{findByTokenHashResult: session}
		uc := NewAuthenticateSessionUseCase(sessions, &fakeSessionTokenGenerator{}, (&fakeClock{now: now}).Now, authSessionIdleTTL)

		if _, err := uc.Execute(context.Background(), "raw-token"); !errors.Is(err, domain.ErrInactiveAdministratorSession) {
			t.Fatalf("expected ErrInactiveAdministratorSession, got %v", err)
		}
		if sessions.updateIdleCalled {
			t.Fatal("no extension should happen for an expired session")
		}
	})

	t.Run("revoked", func(t *testing.T) {
		t.Parallel()
		session := newTestSession(t, now)
		if err := session.Revoke(now.Add(-time.Minute)); err != nil {
			t.Fatalf("revoke: %v", err)
		}
		sessions := &fakeAdministratorSessionRepository{findByTokenHashResult: session}
		uc := NewAuthenticateSessionUseCase(sessions, &fakeSessionTokenGenerator{}, (&fakeClock{now: now}).Now, authSessionIdleTTL)

		if _, err := uc.Execute(context.Background(), "raw-token"); !errors.Is(err, domain.ErrInactiveAdministratorSession) {
			t.Fatalf("expected ErrInactiveAdministratorSession, got %v", err)
		}
		if sessions.updateIdleCalled {
			t.Fatal("no extension should happen for a revoked session")
		}
	})
}

func TestAuthenticateSessionInfrastructureErrorsPropagate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	t.Run("find by token hash fails", func(t *testing.T) {
		t.Parallel()
		wantErr := errors.New("database unavailable")
		sessions := &fakeAdministratorSessionRepository{findByTokenHashErr: wantErr}
		uc := NewAuthenticateSessionUseCase(sessions, &fakeSessionTokenGenerator{}, (&fakeClock{now: now}).Now, authSessionIdleTTL)

		if _, err := uc.Execute(context.Background(), "raw-token"); !errors.Is(err, wantErr) {
			t.Fatalf("expected %v, got %v", wantErr, err)
		}
	})

	t.Run("update idle expiry fails", func(t *testing.T) {
		t.Parallel()
		wantErr := errors.New("database unavailable")
		session := newTestSession(t, now)
		sessions := &fakeAdministratorSessionRepository{findByTokenHashResult: session, updateIdleErr: wantErr}
		uc := NewAuthenticateSessionUseCase(sessions, &fakeSessionTokenGenerator{}, (&fakeClock{now: now}).Now, authSessionIdleTTL)

		if _, err := uc.Execute(context.Background(), "raw-token"); !errors.Is(err, wantErr) {
			t.Fatalf("expected %v, got %v", wantErr, err)
		}
	})
}
