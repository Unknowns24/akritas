package administratorsession_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func seedAdministrator(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()
	id := uuid.New()
	now := time.Now().UTC()
	if err := db.Create(&model.Administrator{
		ID: id, Email: uuid.NewString() + "@example.com", DisplayName: "Akritas Administrator",
		PasswordHash: "hash", EncryptedTOTPSecret: []byte("cipher"), CreatedAt: now, UpdatedAt: now,
	}).Error; err != nil {
		t.Fatalf("seed administrator: %v", err)
	}
	return id
}
