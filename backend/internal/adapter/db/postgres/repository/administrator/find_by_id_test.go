package administrator_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
)

func TestFindByIDReturnsCreatedAdministrator(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	admin := newDomainAdministrator(t, "admin@example.com")
	if err := repo.Create(context.Background(), admin, "hash"); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repo.FindByID(context.Background(), admin.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find the created administrator")
	}
	if found.ID != admin.ID || found.Email != admin.Email || found.DisplayName != admin.DisplayName {
		t.Fatalf("found administrator does not match created one: %+v", found)
	}
}

func TestFindByIDReturnsNilForMissingID(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	found, err := repo.FindByID(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatalf("expected nil for a missing id, got %+v", found)
	}
}
