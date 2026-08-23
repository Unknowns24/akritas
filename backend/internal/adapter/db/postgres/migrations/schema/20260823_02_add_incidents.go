package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_02_AddIncidents() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_02_add_incidents",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("CREATE SEQUENCE incident_key_sequence START 1").Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE TABLE incidents (
                id uuid PRIMARY KEY,
                key varchar(64) NOT NULL UNIQUE,
                project_id uuid NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
                fingerprint varchar(255) NOT NULL,
                severity varchar(16) NOT NULL CHECK (severity IN ('critical','error','warning','info')),
                phase varchar(32) NOT NULL CHECK (phase IN ('detected','investigating','publishing_issue','remediating','completed','failed')),
                terminal_outcome varchar(64) NULL,
                first_seen_at timestamptz NOT NULL,
                last_seen_at timestamptz NOT NULL,
                occurrence_count bigint NOT NULL CHECK (occurrence_count > 0),
                title varchar(500) NOT NULL,
                summary varchar(5000) NOT NULL DEFAULT '',
                root_cause_status varchar(32) NULL,
                resolution_status varchar(32) NULL,
                confidence double precision NULL CHECK (confidence BETWEEN 0 AND 1),
                github_issue_reference jsonb NULL,
                pull_request_reference jsonb NULL,
                CONSTRAINT ck_incident_seen_order CHECK (last_seen_at >= first_seen_at)
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE INDEX ix_incidents_grouping ON incidents(project_id, fingerprint, last_seen_at DESC)",
				"CREATE INDEX ix_incidents_list ON incidents(last_seen_at DESC, id DESC)",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("DROP TABLE IF EXISTS incidents").Error; err != nil {
				return err
			}
			return tx.Exec("DROP SEQUENCE IF EXISTS incident_key_sequence").Error
		},
	}
}
