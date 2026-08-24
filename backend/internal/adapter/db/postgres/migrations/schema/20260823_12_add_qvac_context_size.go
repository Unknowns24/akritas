package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260823_12_AddQvacContextSize() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260823_12_add_qvac_context_size",
		Migrate: func(tx *gorm.DB) error {
			for _, statement := range []string{
				`ALTER TABLE qvac_configurations
				 ADD COLUMN IF NOT EXISTS context_size integer`,
				`ALTER TABLE qvac_configurations
				 ALTER COLUMN context_size SET DEFAULT 16384`,
				`UPDATE qvac_configurations SET context_size = 16384 WHERE context_size IS NULL OR context_size = 0`,
				`ALTER TABLE qvac_configurations
				 ALTER COLUMN context_size SET NOT NULL`,
				`ALTER TABLE qvac_configurations
				 DROP CONSTRAINT IF EXISTS ck_qvac_context_size`,
				`ALTER TABLE qvac_configurations
				 ADD CONSTRAINT ck_qvac_context_size CHECK (context_size BETWEEN 4096 AND 131072)`,
			} {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE qvac_configurations DROP COLUMN IF EXISTS context_size`).Error
		},
	}
}
