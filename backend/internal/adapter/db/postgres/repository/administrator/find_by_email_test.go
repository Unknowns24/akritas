package administrator_test

import (
	"context"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
)

func TestFindByEmailReturnsCredentials(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	admin := newDomainAdministrator(t, "admin@example.com")
	passwordHash := "$argon2id$v=19$m=19456,t=2,p=1$salt$hash"
	if err := repo.Create(context.Background(), admin, passwordHash); err != nil {
		t.Fatalf("create: %v", err)
	}

	creds, err := repo.FindByEmail(context.Background(), "admin@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds == nil {
		t.Fatal("expected to find credentials for the created administrator")
	}
	if creds.Administrator.ID != admin.ID {
		t.Fatalf("administrator id mismatch: %+v", creds.Administrator)
	}
	if creds.PasswordHash != passwordHash {
		t.Fatalf("password hash mismatch: got %q", creds.PasswordHash)
	}
	if creds.Administrator.LastAcceptedTOTPPeriod != -1 {
		t.Fatalf("LastAcceptedTOTPPeriod = %d, want -1 right after Create", creds.Administrator.LastAcceptedTOTPPeriod)
	}
}

func TestFindByEmailReturnsNilForMissingEmail(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	creds, err := repo.FindByEmail(context.Background(), "nobody@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds != nil {
		t.Fatalf("expected nil for a missing email, got %+v", creds)
	}
}
