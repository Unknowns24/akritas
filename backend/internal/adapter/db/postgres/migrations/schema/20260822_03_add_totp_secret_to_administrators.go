package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// administrators is guaranteed empty until this migration's Create call
// runs for the first time, so adding a NOT NULL column here needs no
// backfill.
//
// Migrate is idempotent: model.Administrator is a live Go struct, so on a
// brand-new database migration 20260822_01's AutoMigrate already creates
// this column from the current struct shape, before this migration ever
// runs. On a database that already applied migration 01 back when the
// struct didn't have the field yet (a real pre-existing PB-061 deployment),
// the column is genuinely missing and this migration adds it.
func Migration20260822_03_AddTotpSecretToAdministrators() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_03_add_totp_secret_to_administrators",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(&model.Administrator{}, "EncryptedTOTPSecret") {
				return nil
			}
			return tx.Migrator().AddColumn(&model.Administrator{}, "EncryptedTOTPSecret")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&model.Administrator{}, "EncryptedTOTPSecret")
		},
	}
}
