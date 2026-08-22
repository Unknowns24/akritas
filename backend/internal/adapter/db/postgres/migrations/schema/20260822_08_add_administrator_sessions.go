package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_08_AddAdministratorSessions() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_08_add_administrator_sessions",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE administrator_sessions (
                id uuid PRIMARY KEY,
                administrator_id uuid NOT NULL REFERENCES administrators(id) ON DELETE CASCADE,
                token_hash char(64) NOT NULL UNIQUE,
                authenticated_at timestamptz NOT NULL,
                idle_expires_at timestamptz NOT NULL,
                absolute_expires_at timestamptz NOT NULL,
                revoked_at timestamptz NULL,
                CONSTRAINT ck_administrator_sessions_expiry CHECK (
                    idle_expires_at > authenticated_at AND
                    absolute_expires_at > authenticated_at AND
                    idle_expires_at <= absolute_expires_at
                )
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS administrator_sessions").Error },
	}
}
