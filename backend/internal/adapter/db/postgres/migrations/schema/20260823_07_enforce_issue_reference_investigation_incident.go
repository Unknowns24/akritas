package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_07_EnforceIssueReferenceInvestigationIncident() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_07_enforce_issue_reference_investigation_incident",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`ALTER TABLE investigations
                ADD CONSTRAINT ux_investigations_id_incident_id UNIQUE (id, incident_id)`).Error; err != nil {
				return err
			}
			return tx.Exec(`ALTER TABLE github_issue_references
                ADD CONSTRAINT fk_github_issue_references_investigation_incident
                FOREIGN KEY (investigation_id, incident_id)
                REFERENCES investigations(id, incident_id)
                ON DELETE RESTRICT`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`ALTER TABLE github_issue_references
                DROP CONSTRAINT IF EXISTS fk_github_issue_references_investigation_incident`).Error; err != nil {
				return err
			}
			return tx.Exec(`ALTER TABLE investigations
                DROP CONSTRAINT IF EXISTS ux_investigations_id_incident_id`).Error
		},
	}
}
