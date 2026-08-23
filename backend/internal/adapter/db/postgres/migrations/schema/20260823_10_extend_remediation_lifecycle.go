package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_10_ExtendRemediationLifecycle() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_10_extend_remediation_lifecycle",
		Migrate: func(tx *gorm.DB) error {
			for _, statement := range []string{
				"ALTER TABLE remediations ADD COLUMN investigation_id uuid NULL REFERENCES investigations(id) ON DELETE RESTRICT",
				"ALTER TABLE remediations ADD COLUMN pull_request_number integer NOT NULL DEFAULT 0",
				"ALTER TABLE remediations ADD COLUMN pull_request_url text NOT NULL DEFAULT ''",
				"ALTER TABLE remediations ADD COLUMN pull_request_repository varchar(255) NOT NULL DEFAULT ''",
				"ALTER TABLE remediations ADD COLUMN pull_request_branch varchar(255) NOT NULL DEFAULT ''",
				"ALTER TABLE remediations ADD COLUMN pull_request_created_at timestamptz NULL",
				"ALTER TABLE remediations ADD CONSTRAINT remediations_pull_request_complete_check CHECK ((status = 'pull_request_created') = (pull_request_number > 0 AND pull_request_url <> '' AND pull_request_repository <> '' AND pull_request_branch <> '' AND pull_request_created_at IS NOT NULL))",
				"CREATE UNIQUE INDEX uq_remediations_investigation_id ON remediations (investigation_id) WHERE investigation_id IS NOT NULL",
				"CREATE UNIQUE INDEX uq_remediations_pull_request_head ON remediations (pull_request_repository, pull_request_branch) WHERE pull_request_repository <> '' AND pull_request_branch <> ''",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, statement := range []string{
				"DROP INDEX IF EXISTS uq_remediations_pull_request_head",
				"DROP INDEX IF EXISTS uq_remediations_investigation_id",
				"ALTER TABLE remediations DROP CONSTRAINT IF EXISTS remediations_pull_request_complete_check",
				"ALTER TABLE remediations DROP COLUMN IF EXISTS pull_request_created_at",
				"ALTER TABLE remediations DROP COLUMN IF EXISTS pull_request_branch",
				"ALTER TABLE remediations DROP COLUMN IF EXISTS pull_request_repository",
				"ALTER TABLE remediations DROP COLUMN IF EXISTS pull_request_url",
				"ALTER TABLE remediations DROP COLUMN IF EXISTS pull_request_number",
				"ALTER TABLE remediations DROP COLUMN IF EXISTS investigation_id",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
