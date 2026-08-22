package pendingenrollment_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	pendingenrollment "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func newEnrollment(t *testing.T, email string, createdAt time.Time) *domain.PendingEnrollment {
	t.Helper()
	enrollment, err := domain.NewPendingEnrollment(
		uuid.New(), email, "Akritas Administrator", "$argon2id$v=19$m=19456,t=2,p=1$salt$hash",
		[]byte("ciphertext"), createdAt, createdAt.Add(10*time.Minute),
	)
	if err != nil {
		t.Fatalf("build pending enrollment: %v", err)
	}
	return enrollment
}

func TestSavePersistsPendingEnrollment(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	now := time.Now().UTC()
	enrollment := newEnrollment(t, "admin@example.com", now)

	if err := repo.Save(context.Background(), enrollment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var records []model.PendingEnrollment
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

func TestSaveReplacesPreviousPendingEnrollment(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)

	now := time.Now().UTC()
	first := newEnrollment(t, "first-attempt@example.com", now)
	second := newEnrollment(t, "admin@example.com", now.Add(time.Minute))

	if err := repo.Save(context.Background(), first); err != nil {
		t.Fatalf("unexpected error saving first enrollment: %v", err)
	}
	if err := repo.Save(context.Background(), second); err != nil {
		t.Fatalf("unexpected error saving second enrollment: %v", err)
	}

	var records []model.PendingEnrollment
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
