package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_04_LinkInvestigationHistory() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_04_link_investigation_history",
		Migrate: func(tx *gorm.DB) error {
			statements := []string{
				`ALTER TABLE investigations ADD CONSTRAINT fk_investigations_incident FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE RESTRICT`,
				`ALTER TABLE evidence DROP CONSTRAINT evidence_investigation_id_fkey`,
				`ALTER TABLE evidence ADD CONSTRAINT fk_evidence_investigation FOREIGN KEY (investigation_id) REFERENCES investigations(id) ON DELETE RESTRICT`,
				`DROP INDEX ix_investigations_incident_active`,
				`CREATE UNIQUE INDEX ux_investigations_one_active_per_incident ON investigations (incident_id) WHERE status IN ('pending','running')`,
			}
			for _, statement := range statements {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			statements := []string{
				`DROP INDEX IF EXISTS ux_investigations_one_active_per_incident`,
				`CREATE INDEX ix_investigations_incident_active ON investigations (incident_id) WHERE status IN ('pending','running')`,
				`ALTER TABLE evidence DROP CONSTRAINT IF EXISTS fk_evidence_investigation`,
				`ALTER TABLE evidence ADD CONSTRAINT evidence_investigation_id_fkey FOREIGN KEY (investigation_id) REFERENCES investigations(id) ON DELETE CASCADE`,
				`ALTER TABLE investigations DROP CONSTRAINT IF EXISTS fk_investigations_incident`,
			}
			for _, statement := range statements {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
