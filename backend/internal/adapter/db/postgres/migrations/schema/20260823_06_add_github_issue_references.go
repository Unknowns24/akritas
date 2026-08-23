package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_06_AddGitHubIssueReferences() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_06_add_github_issue_references",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE github_issue_references (
                incident_id uuid NOT NULL REFERENCES incidents(id) ON DELETE RESTRICT,
                investigation_id uuid PRIMARY KEY REFERENCES investigations(id) ON DELETE RESTRICT,
                issue_number integer NOT NULL CHECK (issue_number > 0),
                issue_url text NOT NULL,
                repository varchar(255) NOT NULL,
                created_at timestamptz NOT NULL,
                CONSTRAINT ux_github_issue_references_repository_number UNIQUE (repository, issue_number)
            )`).Error; err != nil {
				return err
			}
			return tx.Exec(`CREATE INDEX ix_github_issue_references_incident_created_at ON github_issue_references (incident_id, created_at DESC)`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("DROP TABLE IF EXISTS github_issue_references").Error
		},
	}
}
