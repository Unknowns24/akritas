package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_12_AddEvidence() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_12_add_evidence",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE evidence (
                id uuid PRIMARY KEY,
                investigation_id uuid NOT NULL REFERENCES investigations(id) ON DELETE CASCADE,
                type varchar(32) NOT NULL CHECK (type IN ('log_excerpt','stack_trace','code_location','commit','diff','validation_result','deployment_metadata')),
                summary varchar(5000) NOT NULL,
                content varchar(50000) NOT NULL DEFAULT '',
                file_path varchar(4096) NOT NULL DEFAULT '',
                line_start integer NULL CHECK (line_start IS NULL OR line_start >= 1),
                line_end integer NULL CHECK (line_end IS NULL OR line_end >= 1),
                commit_sha varchar(64) NOT NULL DEFAULT '',
                patch varchar(100000) NOT NULL DEFAULT '',
                occurred_at timestamptz NULL,
                redacted boolean NOT NULL DEFAULT true CHECK (redacted = true),
                created_at timestamptz NOT NULL,
                CONSTRAINT ck_evidence_lines CHECK ((line_start IS NULL) = (line_end IS NULL))
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE INDEX ix_evidence_investigation_id ON evidence (investigation_id)",
				"CREATE INDEX ix_evidence_created_at_id ON evidence (created_at, id)",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS evidence").Error },
	}
}
