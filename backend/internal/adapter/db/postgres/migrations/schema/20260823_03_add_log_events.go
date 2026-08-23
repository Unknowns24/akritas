package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_03_AddLogEvents() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_03_add_log_events",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE log_events (
                id uuid PRIMARY KEY,
                incident_id uuid NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
                project_id uuid NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
                source_application_id varchar(255) NOT NULL,
                source_instance_id varchar(255) NOT NULL,
                occurrence_key varchar(255) NOT NULL,
                timestamp timestamptz NOT NULL,
                severity varchar(16) NOT NULL CHECK (severity IN ('critical','error','warning','info')),
                message varchar(20000) NOT NULL,
                fingerprint varchar(255) NOT NULL,
                detection_rules jsonb NOT NULL,
                context_before jsonb NOT NULL DEFAULT '[]'::jsonb,
                context_after jsonb NOT NULL DEFAULT '[]'::jsonb,
                raw_context_redacted boolean NOT NULL DEFAULT true,
                CONSTRAINT uq_log_events_occurrence UNIQUE (project_id, occurrence_key)
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE INDEX ix_log_events_incident_timestamp ON log_events(incident_id, timestamp DESC, id DESC)",
				"CREATE INDEX ix_log_events_project_fingerprint ON log_events(project_id, fingerprint)",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS log_events").Error },
	}
}
