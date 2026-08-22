package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func Migration20260822_01_CreateAdministrators() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_01_create_administrators",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Administrator{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&model.Administrator{})
		},
	}
}
