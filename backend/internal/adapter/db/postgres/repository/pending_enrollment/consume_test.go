package pendingenrollment_test

import (
	"context"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	pendingenrollment "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
)

func TestConsumeReturnsEnrollmentExactlyOnce(t *testing.T) {
	db := dbtest.Connect(t)
	repo := pendingenrollment.NewRepository(db)
	enrollment := newEnrollment(t, "admin@example.com", time.Now().UTC())
	if _, err := repo.Replace(context.Background(), enrollment, "new-password-hash"); err != nil {
		t.Fatal(err)
	}
	first, err := repo.Consume(context.Background(), enrollment.ID)
	if err != nil || first == nil || first.PasswordHash != "new-password-hash" {
		t.Fatalf("first consume=%+v err=%v", first, err)
	}
	second, err := repo.Consume(context.Background(), enrollment.ID)
	if err != nil || second != nil {
		t.Fatalf("second consume=%+v err=%v", second, err)
	}
}
