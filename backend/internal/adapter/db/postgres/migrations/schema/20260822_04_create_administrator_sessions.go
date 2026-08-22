package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func Migration20260822_04_CreateAdministratorSessions() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_04_create_administrator_sessions",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.AdministratorSession{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&model.AdministratorSession{})
		},
	}
}
