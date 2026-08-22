package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_07_AddPendingEnrollments() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_07_add_pending_enrollments",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE pending_enrollments (
                id uuid PRIMARY KEY,
                email varchar(254) NOT NULL,
                display_name varchar(100) NOT NULL,
                password_hash text NOT NULL,
                created_at timestamptz NOT NULL,
                expires_at timestamptz NOT NULL CHECK (expires_at > created_at)
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS pending_enrollments").Error },
	}
}
