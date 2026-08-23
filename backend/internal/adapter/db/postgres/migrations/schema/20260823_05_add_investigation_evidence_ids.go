package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_05_AddInvestigationEvidenceIDs() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_05_add_investigation_evidence_ids",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE investigations ADD COLUMN evidence_ids jsonb NOT NULL DEFAULT '[]'::jsonb`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE investigations DROP COLUMN IF EXISTS evidence_ids`).Error
		},
	}
}
