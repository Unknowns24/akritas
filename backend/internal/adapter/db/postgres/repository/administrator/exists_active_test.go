package administrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
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

	now := time.Now().UTC()
	record := &model.Administrator{
		ID:           uuid.New(),
		Email:        "admin@example.com",
		DisplayName:  "Akritas Administrator",
		PasswordHash: "$argon2id$v=19$m=19456,t=2,p=1$salt$hash",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.Create(record).Error; err != nil {
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
