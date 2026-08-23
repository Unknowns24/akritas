package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_05_AddGitHubAppBindings() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_05_add_github_app_bindings",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE github_app_bindings (
                github_account_id uuid PRIMARY KEY REFERENCES github_accounts(id) ON DELETE CASCADE,
                app_id bigint NOT NULL,
                installation_id bigint NOT NULL UNIQUE,
                app_slug varchar(255) NOT NULL,
                client_id varchar(255) NOT NULL,
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL
            )`).Error
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS github_app_bindings").Error },
	}
}
