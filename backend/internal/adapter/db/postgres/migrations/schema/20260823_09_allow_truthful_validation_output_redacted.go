package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_09_AllowTruthfulValidationOutputRedacted() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_09_allow_truthful_validation_output_redacted",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE validation_results
                DROP CONSTRAINT IF EXISTS validation_results_output_redacted_check`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`UPDATE validation_results SET output_redacted = true WHERE output_redacted = false`).Error; err != nil {
				return err
			}
			return tx.Exec(`ALTER TABLE validation_results
                ADD CONSTRAINT validation_results_output_redacted_check CHECK (output_redacted = true)`).Error
		},
	}
}
