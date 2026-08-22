package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_04_AddGitHubAppRegistrations() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_04_add_github_app_registrations",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE github_app_registrations (
                id uuid PRIMARY KEY,
                display_name varchar(200) NOT NULL,
                account_type varchar(32) NOT NULL CHECK (account_type IN ('personal','organization')),
                account_identifier varchar(255) NOT NULL,
                conversion_state_digest bytea NOT NULL UNIQUE,
                installation_state_digest bytea NULL UNIQUE,
                status varchar(32) NOT NULL CHECK (status IN ('created','converted','completed')),
                app_id bigint NULL,
                app_slug varchar(255) NOT NULL DEFAULT '',
                app_name varchar(255) NOT NULL DEFAULT '',
                client_id varchar(255) NOT NULL DEFAULT '',
                expires_at timestamptz NOT NULL,
                conversion_consumed_at timestamptz NULL,
                installation_consumed_at timestamptz NULL,
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS github_app_registrations").Error },
	}
}
