package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// investigations.incident_id has no foreign key yet: H2 (Detection +
// Incidents) has not merged a real `incidents` table into this repository.
// A follow-up migration adds the FK once H2 lands.
func SCHEMA_20260822_10_AddInvestigations() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_10_add_investigations",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`CREATE TABLE investigations (
                id uuid PRIMARY KEY,
                incident_id uuid NOT NULL,
                status varchar(32) NOT NULL CHECK (status IN ('pending','running','completed','failed')),
                created_at timestamptz NOT NULL,
                started_at timestamptz NULL,
                finished_at timestamptz NULL,
                summary varchar(10000) NOT NULL DEFAULT '',
                root_cause varchar(20000) NOT NULL DEFAULT '',
                root_cause_status varchar(32) NULL CHECK (root_cause_status IN ('identified','suspected','unknown')),
                resolution_status varchar(32) NULL CHECK (resolution_status IN ('fixable','requires_human')),
                confidence double precision NULL CHECK (confidence IS NULL OR (confidence >= 0 AND confidence <= 1)),
                hypotheses jsonb NOT NULL DEFAULT '[]'::jsonb,
                relevant_files jsonb NOT NULL DEFAULT '[]'::jsonb,
                relevant_commits jsonb NOT NULL DEFAULT '[]'::jsonb,
                recommended_actions jsonb NOT NULL DEFAULT '[]'::jsonb,
                evidence_count integer NOT NULL DEFAULT 0 CHECK (evidence_count >= 0),
                failure_user_message varchar(1000) NOT NULL DEFAULT ''
            )`).Error; err != nil {
				return err
			}
			for _, statement := range []string{
				"CREATE INDEX ix_investigations_incident_id ON investigations (incident_id)",
				"CREATE INDEX ix_investigations_incident_active ON investigations (incident_id) WHERE status IN ('pending','running')",
				"CREATE INDEX ix_investigations_created_at_id ON investigations (created_at, id)",
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error { return tx.Exec("DROP TABLE IF EXISTS investigations").Error },
	}
}
