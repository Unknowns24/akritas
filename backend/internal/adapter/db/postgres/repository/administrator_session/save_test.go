package administratorsession_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	administratorsession "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator_session"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestSavePersistsAdministratorSession(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administratorsession.NewRepository(db)

	administratorID := uuid.New()
	now := time.Now().UTC()
	if err := db.Create(&model.Administrator{
		ID: administratorID, Email: "admin@example.com", DisplayName: "Akritas Administrator",
		PasswordHash: "hash", EncryptedTOTPSecret: []byte("cipher"), CreatedAt: now, UpdatedAt: now,
	}).Error; err != nil {
		t.Fatalf("seed administrator: %v", err)
	}

	session, err := domain.NewAdministratorSession(uuid.New(), administratorID, now, now.Add(12*time.Hour), now.Add(168*time.Hour))
	if err != nil {
		t.Fatalf("build session: %v", err)
	}
	tokenHash := "deadbeef"

	if err := repo.Save(context.Background(), session, tokenHash); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var record model.AdministratorSession
	if err := db.First(&record, "id = ?", session.ID).Error; err != nil {
		t.Fatalf("query session: %v", err)
	}
	if record.AdministratorID != administratorID || record.TokenHash != tokenHash {
		t.Fatalf("persisted row does not match input: %+v", record)
	}
}
