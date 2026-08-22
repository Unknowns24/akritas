package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func Migration20260822_02_CreatePendingEnrollments() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_02_create_pending_enrollments",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.PendingEnrollment{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&model.PendingEnrollment{})
		},
	}
}
