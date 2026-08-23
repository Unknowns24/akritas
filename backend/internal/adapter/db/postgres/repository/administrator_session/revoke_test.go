package administratorsession_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	administratorsession "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator_session"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestRevokePersistsRevokedAt(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administratorsession.NewRepository(db)
	administratorID := seedAdministrator(t, db)

	now := time.Now().UTC()
	session, err := domain.NewAdministratorSession(uuid.New(), administratorID, now, now.Add(12*time.Hour), now.Add(168*time.Hour))
	if err != nil {
		t.Fatalf("build session: %v", err)
	}
	if err := repo.Save(context.Background(), session, "deadbeef"); err != nil {
		t.Fatalf("save: %v", err)
	}

	revokedAt := now.Add(5 * time.Minute)
	if err := repo.Revoke(context.Background(), session.ID, revokedAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByTokenHash(context.Background(), "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find the session")
	}
	if found.RevokedAt == nil || !found.RevokedAt.Equal(revokedAt) {
		t.Fatalf("RevokedAt = %v, want %v", found.RevokedAt, revokedAt)
	}
	if found.IsActive(revokedAt.Add(time.Second)) {
		t.Fatal("a revoked session must not be active")
	}
}
