package administratorsession_test

import (
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"

	administratorrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func seedAdministrator(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()
	now := time.Now().UTC()
	administrator, err := domain.NewAdministrator(uuid.New(), uuid.NewString()+"@example.com", "Akritas Administrator", now)
	if err != nil {
		t.Fatalf("build administrator: %v", err)
	}
	if err := administratorrepo.NewRepository(db).Create(t.Context(), administrator, "hash"); err != nil {
		t.Fatalf("seed administrator: %v", err)
	}
	return administrator.ID
}
