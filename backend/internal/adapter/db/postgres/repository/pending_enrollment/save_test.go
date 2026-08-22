package pendingenrollment_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	pendingenrollment "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func newEnrollment(t *testing.T, email string, createdAt time.Time) *domain.PendingEnrollment {
	t.Helper()
	enrollment, err := domain.NewPendingEnrollment(
		uuid.New(), email, "Akritas Administrator", createdAt, createdAt.Add(10*time.Minute),
	)
	if err != nil {
		t.Fatalf("build pending enrollment: %v", err)
	}
	return enrollment
}

func TestReplacePersistsPendingEnrollment(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	now := time.Now().UTC()
	enrollment := newEnrollment(t, "admin@example.com", now)

	if _, err := repo.Replace(context.Background(), enrollment, "password-hash"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var records []domain.PendingEnrollment
	if err := db.Find(&records).Error; err != nil {
		t.Fatalf("query pending enrollments: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected exactly 1 row, got %d", len(records))
	}
	if records[0].ID != enrollment.ID || records[0].Email != enrollment.Email {
		t.Fatalf("persisted row does not match saved enrollment: %+v", records[0])
	}
}

func TestReplaceReturnsAndReplacesPreviousPendingEnrollment(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	now := time.Now().UTC()
	first := newEnrollment(t, "first-attempt@example.com", now)
	second := newEnrollment(t, "admin@example.com", now.Add(time.Minute))

	if previous, err := repo.Replace(context.Background(), first, "hash-one"); err != nil || previous != nil {
		t.Fatalf("unexpected error saving first enrollment: %v", err)
	}
	previous, err := repo.Replace(context.Background(), second, "hash-two")
	if err != nil {
		t.Fatalf("unexpected error saving second enrollment: %v", err)
	}
	if previous == nil || *previous != first.ID {
		t.Fatalf("expected replaced enrollment id %s, got %v", first.ID, previous)
	}

	var records []domain.PendingEnrollment
	if err := db.Find(&records).Error; err != nil {
		t.Fatalf("query pending enrollments: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected exactly 1 surviving row after two saves, got %d", len(records))
	}
	if records[0].ID != second.ID || records[0].Email != second.Email {
		t.Fatalf("the surviving row must be the most recent enrollment, got %+v", records[0])
	}
}
