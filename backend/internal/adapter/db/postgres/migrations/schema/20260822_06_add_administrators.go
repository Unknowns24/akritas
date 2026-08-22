package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_06_AddAdministrators() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_06_add_administrators",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE administrators (
                id uuid PRIMARY KEY,
                email varchar(254) NOT NULL UNIQUE,
                display_name varchar(100) NOT NULL,
                password_hash text NOT NULL,
                last_accepted_totp_period bigint NOT NULL DEFAULT -1 CHECK (last_accepted_totp_period >= -1),
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS administrators").Error },
	}
}
