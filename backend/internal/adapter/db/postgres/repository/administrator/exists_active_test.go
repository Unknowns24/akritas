package administrator_test

import (
	"context"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	"testing"
)

func TestExistsActive(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	exists, err := repo.ExistsActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatal("expected no administrator in an empty table")
	}

	record := newDomainAdministrator(t, "admin@example.com")
	if err := repo.Create(context.Background(), record, "hash"); err != nil {
		t.Fatalf("seed administrator: %v", err)
	}

	exists, err = repo.ExistsActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("expected an administrator to exist after insert")
	}
}
