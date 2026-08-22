package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_01_AddGitHubAccounts() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_01_add_github_accounts",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE github_accounts (
                id uuid PRIMARY KEY,
                display_name varchar(200) NOT NULL,
                account_type varchar(32) NOT NULL CHECK (account_type IN ('personal','organization')),
                authentication_method varchar(32) NOT NULL CHECK (authentication_method IN ('personal_access_token','github_app')),
                account_identifier varchar(255) NOT NULL,
                authentication_status varchar(32) NOT NULL CHECK (authentication_status IN ('pending','connected','authentication_failed','unavailable')),
                credential_configured boolean NOT NULL DEFAULT false,
                repository_count integer NOT NULL DEFAULT 0 CHECK (repository_count >= 0),
                last_checked_at timestamptz NULL,
                manage_url text NOT NULL DEFAULT '',
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL,
                CONSTRAINT uq_github_accounts_identity UNIQUE (authentication_method, account_identifier)
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS github_accounts").Error },
	}
}
