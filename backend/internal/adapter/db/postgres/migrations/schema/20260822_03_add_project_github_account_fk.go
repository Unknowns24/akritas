package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_03_AddProjectGitHubAccountFK() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_03_add_project_github_account_fk",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE projects
				ADD CONSTRAINT fk_projects_github_account
				FOREIGN KEY (github_account_id) REFERENCES github_accounts(id)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE projects DROP CONSTRAINT IF EXISTS fk_projects_github_account`).Error
		},
	}
}
