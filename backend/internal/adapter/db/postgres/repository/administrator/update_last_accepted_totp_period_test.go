package administrator_test

import (
	"context"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
)

func TestConsumeTOTPPeriodIsStrictCompareAndSet(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	admin := newDomainAdministrator(t, "admin@example.com")
	if err := repo.Create(context.Background(), admin, "hash"); err != nil {
		t.Fatalf("create: %v", err)
	}

	consumed, err := repo.ConsumeTOTPPeriod(context.Background(), admin.ID, "hash", 123456789)
	if err != nil || !consumed {
		t.Fatalf("unexpected error: %v", err)
	}
	consumed, err = repo.ConsumeTOTPPeriod(context.Background(), admin.ID, "hash", 123456789)
	if err != nil || consumed {
		t.Fatalf("same period must not be consumed twice: consumed=%v err=%v", consumed, err)
	}
	consumed, err = repo.ConsumeTOTPPeriod(context.Background(), admin.ID, "stale-hash", 123456790)
	if err != nil || consumed {
		t.Fatalf("stale password generation must not consume period: consumed=%v err=%v", consumed, err)
	}

	creds, err := repo.FindByEmail(context.Background(), "admin@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds == nil {
		t.Fatal("expected to find the administrator")
	}
	if creds.Administrator.LastAcceptedTOTPPeriod != 123456789 {
		t.Fatalf("LastAcceptedTOTPPeriod = %d, want 123456789", creds.Administrator.LastAcceptedTOTPPeriod)
	}

	// UpdateColumn must bypass GORM's auto-updated_at convention: the
	// domain reconstruction always sets UpdatedAt == CreatedAt, and this
	// must stay true (updated_at must not silently drift in the DB).
	if !creds.Administrator.UpdatedAt.Equal(creds.Administrator.CreatedAt) {
		t.Fatalf("UpdatedAt (%v) drifted from CreatedAt (%v) after updating an unrelated column",
			creds.Administrator.UpdatedAt, creds.Administrator.CreatedAt)
	}
}
