package administrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
)

func TestRotateCredentialsRejectsStalePasswordHash(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)
	admin := newDomainAdministrator(t, "admin@example.com")
	if err := repo.Create(context.Background(), admin, "old-hash"); err != nil {
		t.Fatal(err)
	}
	if rotated, err := repo.RotateCredentials(context.Background(), admin.ID, "stale-hash", "new-hash", 42, time.Now().UTC()); err != nil || rotated {
		t.Fatalf("stale rotate=%v err=%v", rotated, err)
	}
	updatedAt := time.Now().UTC().Add(time.Second)
	if rotated, err := repo.RotateCredentials(context.Background(), admin.ID, "old-hash", "new-hash", 42, updatedAt); err != nil || !rotated {
		t.Fatalf("rotate=%v err=%v", rotated, err)
	}
	creds, err := repo.FindByEmail(context.Background(), admin.Email)
	if err != nil || creds == nil {
		t.Fatalf("find=%+v err=%v", creds, err)
	}
	if creds.PasswordHash != "new-hash" || creds.Administrator.LastAcceptedTOTPPeriod != 42 || !creds.Administrator.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("credentials not rotated: %+v", creds)
	}
}
