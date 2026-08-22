package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAdministratorSessionLifecycle(t *testing.T) {
	t.Parallel()

	authenticatedAt := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session, err := NewAdministratorSession(
		uuid.New(),
		uuid.New(),
		authenticatedAt,
		authenticatedAt.Add(12*time.Hour),
		authenticatedAt.Add(7*24*time.Hour),
	)
	if err != nil {
		t.Fatalf("valid session rejected: %v", err)
	}
	if !session.IsActive(authenticatedAt.Add(time.Hour)) {
		t.Fatal("session should be active before expiration")
	}
	if session.IsActive(session.IdleExpiresAt) || session.IsActive(session.AbsoluteExpiresAt) {
		t.Fatal("expiration boundary must be exclusive")
	}

	revokedAt := authenticatedAt.Add(2 * time.Hour)
	if err := session.Revoke(revokedAt); err != nil {
		t.Fatalf("revoke failed: %v", err)
	}
	if err := session.Revoke(revokedAt.Add(time.Minute)); err != nil {
		t.Fatalf("revoke must be idempotent: %v", err)
	}
	if session.IsActive(revokedAt.Add(time.Second)) {
		t.Fatal("revoked session must not be active")
	}
}

func TestAdministratorSessionRejectsInvalidTimes(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	_, err := NewAdministratorSession(uuid.New(), uuid.New(), now, now.Add(8*time.Hour), now.Add(4*time.Hour))
	if !errors.Is(err, ErrInvalidAdministratorSession) {
		t.Fatalf("expected invalid session error, got %v", err)
	}

	session, err := NewAdministratorSession(uuid.New(), uuid.New(), now, now.Add(time.Hour), now.Add(2*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Revoke(now.Add(-time.Second)); !errors.Is(err, ErrAdministratorSessionTransition) {
		t.Fatalf("expected transition error, got %v", err)
	}
}

func TestAdministratorSessionExtendIdle(t *testing.T) {
	t.Parallel()

	authenticatedAt := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session, err := NewAdministratorSession(
		uuid.New(), uuid.New(), authenticatedAt,
		authenticatedAt.Add(time.Hour), authenticatedAt.Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("valid session rejected: %v", err)
	}

	now := authenticatedAt.Add(30 * time.Minute)
	if err := session.ExtendIdle(now, time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !session.IdleExpiresAt.Equal(now.Add(time.Hour)) {
		t.Fatalf("IdleExpiresAt = %v, want %v", session.IdleExpiresAt, now.Add(time.Hour))
	}
}

func TestAdministratorSessionExtendIdleCapsAtAbsoluteDeadline(t *testing.T) {
	t.Parallel()

	authenticatedAt := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session, err := NewAdministratorSession(
		uuid.New(), uuid.New(), authenticatedAt,
		authenticatedAt.Add(10*time.Minute), authenticatedAt.Add(20*time.Minute),
	)
	if err != nil {
		t.Fatalf("valid session rejected: %v", err)
	}

	// now is still within the current idle window, but idleTTL alone would
	// push the new idle deadline well past the absolute one.
	now := authenticatedAt.Add(5 * time.Minute)
	if err := session.ExtendIdle(now, time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !session.IdleExpiresAt.Equal(session.AbsoluteExpiresAt) {
		t.Fatalf("IdleExpiresAt = %v, want capped at AbsoluteExpiresAt %v", session.IdleExpiresAt, session.AbsoluteExpiresAt)
	}
}

func TestAdministratorSessionExtendIdleRejectsInactiveSession(t *testing.T) {
	t.Parallel()

	authenticatedAt := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	session, err := NewAdministratorSession(
		uuid.New(), uuid.New(), authenticatedAt,
		authenticatedAt.Add(time.Hour), authenticatedAt.Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("valid session rejected: %v", err)
	}
	if err := session.Revoke(authenticatedAt.Add(time.Minute)); err != nil {
		t.Fatalf("revoke failed: %v", err)
	}

	before := session.IdleExpiresAt
	err = session.ExtendIdle(authenticatedAt.Add(2*time.Minute), time.Hour)
	if !errors.Is(err, ErrAdministratorSessionTransition) {
		t.Fatalf("expected transition error extending a revoked session, got %v", err)
	}
	if !session.IdleExpiresAt.Equal(before) {
		t.Fatal("a rejected extension must not modify IdleExpiresAt")
	}
}
