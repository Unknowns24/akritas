package schema

import (
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_02_AddProjects() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_02_add_projects",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&domain.Project{}); err != nil {
				return err
			}
			return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS ux_projects_name_lower ON projects (LOWER(name))`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP INDEX IF EXISTS ux_projects_name_lower`).Error; err != nil {
				return err
			}
			if tx.Migrator().HasTable(&domain.Project{}) {
				return tx.Migrator().DropTable(&domain.Project{})
			}
			return nil
		},
	}
}
