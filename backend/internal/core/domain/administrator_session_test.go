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
