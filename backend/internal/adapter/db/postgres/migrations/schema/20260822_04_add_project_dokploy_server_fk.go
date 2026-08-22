package schema

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func SCHEMA_20260822_04_AddProjectDokployServerFK() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260822_04_add_project_dokploy_server_fk",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
				ALTER TABLE projects
				ADD CONSTRAINT fk_projects_dokploy_server
				FOREIGN KEY (dokploy_server_id) REFERENCES dokploy_servers(id)
			`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`ALTER TABLE projects DROP CONSTRAINT IF EXISTS fk_projects_dokploy_server`).Error
		},
	}
}
