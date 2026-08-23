package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_01_AddMonitoringCheckpoints() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_01_add_monitoring_checkpoints",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE monitoring_checkpoints (
                id uuid PRIMARY KEY,
                project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                source_application_id varchar(255) NOT NULL,
                source_instance_id varchar(255) NOT NULL,
                is_current boolean NOT NULL DEFAULT true,
                initial_backfill_pending boolean NOT NULL DEFAULT false,
                cursor_timestamp timestamptz NULL,
                cursor_ordinal integer NOT NULL DEFAULT 0 CHECK (cursor_ordinal >= 0),
                cursor_content_hash varchar(128) NOT NULL DEFAULT '',
                version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
                assembly_state jsonb NOT NULL DEFAULT '{"recent_records":[],"open_records":[],"pending":[]}'::jsonb,
                next_finalize_at timestamptz NULL,
                created_at timestamptz NOT NULL,
                updated_at timestamptz NOT NULL
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE UNIQUE INDEX uq_monitoring_checkpoints_current_project ON monitoring_checkpoints(project_id) WHERE is_current",
				"CREATE INDEX ix_monitoring_checkpoints_finalize ON monitoring_checkpoints(next_finalize_at) WHERE next_finalize_at IS NOT NULL",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS monitoring_checkpoints").Error },
	}
}
