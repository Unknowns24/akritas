package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// administrators is guaranteed empty until this migration's Create call
// runs for the first time, so adding a NOT NULL column here needs no
// backfill.
func Migration20260822_03_AddTotpSecretToAdministrators() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_03_add_totp_secret_to_administrators",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&model.Administrator{}, "EncryptedTOTPSecret")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&model.Administrator{}, "EncryptedTOTPSecret")
		},
	}
}
