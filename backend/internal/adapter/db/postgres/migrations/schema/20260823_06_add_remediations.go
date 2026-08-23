package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// remediations is deliberately minimal (id, incident_id, status, branch
// naming, and terminal-failure text only): full lifecycle persistence
// (Changes, PullRequestReference, status transitions beyond in_progress)
// is deferred to AKR-49/55+. It exists now so validation_results has a
// real foreign key from day one.
func SCHEMA_20260823_06_AddRemediations() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_06_add_remediations",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE remediations (
                id uuid PRIMARY KEY,
                incident_id uuid NOT NULL REFERENCES incidents(id) ON DELETE RESTRICT,
                status varchar(32) NOT NULL CHECK (status IN ('planned','in_progress','validated','failed','pull_request_created')),
                branch_name varchar(255) NOT NULL DEFAULT '',
                changes_summary varchar(2000) NOT NULL DEFAULT '',
                failure_user_message varchar(1000) NOT NULL DEFAULT '',
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE INDEX ix_remediations_incident_id ON remediations (incident_id)",
				"CREATE UNIQUE INDEX uq_remediations_branch_name ON remediations (branch_name) WHERE branch_name <> ''",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS remediations").Error },
	}
}
