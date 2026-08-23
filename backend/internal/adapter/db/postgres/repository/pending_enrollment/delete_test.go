package pendingenrollment_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	pendingenrollment "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
)

func TestDeleteRemovesEnrollment(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	enrollment := newEnrollment(t, "admin@example.com", time.Now().UTC())
	if _, err := repo.Replace(context.Background(), enrollment, "password-hash"); err != nil {
		t.Fatalf("save: %v", err)
	}

	if err := repo.Delete(context.Background(), enrollment.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(context.Background(), enrollment.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatal("expected the enrollment to be gone after Delete")
	}
}

func TestDeleteIsIdempotent(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	if err := repo.Delete(context.Background(), uuid.New()); err != nil {
		t.Fatalf("deleting a non-existent id must not error, got: %v", err)
	}
}
