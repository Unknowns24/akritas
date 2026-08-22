package pendingenrollment_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	pendingenrollment "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
)

func TestFindByIDReturnsSavedEnrollment(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	now := time.Now().UTC()
	enrollment := newEnrollment(t, "admin@example.com", now)

	if _, err := repo.Replace(context.Background(), enrollment, "password-hash"); err != nil {
		t.Fatalf("save: %v", err)
	}

	found, err := repo.FindByID(context.Background(), enrollment.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find the saved enrollment")
	}
	if found.Enrollment.ID != enrollment.ID || found.Enrollment.Email != enrollment.Email || found.PasswordHash != "password-hash" {
		t.Fatalf("found enrollment does not match saved one: %+v", found)
	}
}

func TestFindByIDReturnsNilForMissingID(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	found, err := repo.FindByID(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Fatalf("expected nil for a missing id, got %+v", found)
	}
}
