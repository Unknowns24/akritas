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

func TestFindByTokenHashReturnsSavedSession(t *testing.T) {
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

	found, err := repo.FindByTokenHash(context.Background(), "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find the saved session")
	}
	if found.ID != session.ID || found.AdministratorID != administratorID {
		t.Fatalf("found session does not match saved one: %+v", found)
	}
	if found.RevokedAt != nil {
		t.Fatal("a freshly saved session must not be revoked")
	}
}

func TestFindByTokenHashReturnsNilForMissingHash(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administratorsession.NewRepository(db)

	found, err := repo.FindByTokenHash(context.Background(), "does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatalf("expected nil for a missing hash, got %+v", found)
	}
}
