package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// Migrate is idempotent, same reasoning as migration 20260822_03: on a
// brand-new database, migration 20260822_01's AutoMigrate already creates
// this column from the current (already-extended) model.Administrator
// struct. This migration only does real work on a database that applied
// migration 01 before the struct had the field.
func Migration20260822_05_AddLastAcceptedTotpPeriodToAdministrators() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_05_add_last_accepted_totp_period_to_administrators",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(&model.Administrator{}, "LastAcceptedTOTPPeriod") {
				return nil
			}
			return tx.Migrator().AddColumn(&model.Administrator{}, "LastAcceptedTOTPPeriod")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&model.Administrator{}, "LastAcceptedTOTPPeriod")
		},
	}
}
