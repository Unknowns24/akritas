package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_07_AddValidationResults() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_07_add_validation_results",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE validation_results (
                id uuid PRIMARY KEY,
                remediation_id uuid NOT NULL REFERENCES remediations(id) ON DELETE RESTRICT,
                type varchar(32) NOT NULL CHECK (type IN ('test','build','static_analysis')),
                name varchar(255) NOT NULL,
                status varchar(32) NOT NULL CHECK (status IN ('pending','running','passed','failed')),
                created_at timestamptz NOT NULL,
                started_at timestamptz NULL,
                finished_at timestamptz NULL,
                summary varchar(5000) NOT NULL DEFAULT '',
                output_excerpt varchar(50000) NOT NULL DEFAULT '',
                output_redacted boolean NOT NULL DEFAULT true CHECK (output_redacted = true),
                CONSTRAINT ck_validation_results_times CHECK (
                    (status = 'pending' AND started_at IS NULL AND finished_at IS NULL) OR
                    (status = 'running' AND started_at IS NOT NULL AND finished_at IS NULL) OR
                    (status IN ('passed','failed') AND started_at IS NOT NULL AND finished_at IS NOT NULL)
                )
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE INDEX ix_validation_results_remediation_id ON validation_results (remediation_id)",
				"CREATE INDEX ix_validation_results_created_at_id ON validation_results (created_at, id)",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS validation_results").Error },
	}
}
