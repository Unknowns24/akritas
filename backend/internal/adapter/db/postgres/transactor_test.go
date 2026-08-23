package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestTransactorRollsBackOnError(t *testing.T) {
	db := dbtest.Connect(t)
	transactor := postgres.NewTransactor(db)
	administrators := administrator.NewRepository(db)

	admin, err := domain.NewAdministrator(uuid.New(), "admin@example.com", "Akritas Administrator", time.Now().UTC())
	if err != nil {
		t.Fatalf("build administrator: %v", err)
	}

	failure := errors.New("simulated failure after create")
	err = transactor.WithinTransaction(context.Background(), func(ctx context.Context) error {
		if err := administrators.Create(ctx, admin, "hash"); err != nil {
			return err
		}
		return failure
	})
	if !errors.Is(err, failure) {
		t.Fatalf("expected the simulated failure to propagate, got %v", err)
	}

	exists, err := administrators.ExistsActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatal("the Create inside the failed transaction must have been rolled back")
	}
}

func TestTransactorCommitsOnSuccess(t *testing.T) {
	db := dbtest.Connect(t)
	transactor := postgres.NewTransactor(db)
	administrators := administrator.NewRepository(db)

	admin, err := domain.NewAdministrator(uuid.New(), "admin@example.com", "Akritas Administrator", time.Now().UTC())
	if err != nil {
		t.Fatalf("build administrator: %v", err)
	}

	err = transactor.WithinTransaction(context.Background(), func(ctx context.Context) error {
		return administrators.Create(ctx, admin, "hash")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exists, err := administrators.ExistsActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatal("the Create inside the successful transaction must be visible afterward")
	}
}
