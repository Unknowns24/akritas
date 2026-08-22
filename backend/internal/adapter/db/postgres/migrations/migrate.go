package migrations

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations/schema"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	migrator := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		schema.SCHEMA_20260822_01_AddIntegrationLookups(),
		schema.SCHEMA_20260822_02_AddProjects(),
		schema.SCHEMA_20260822_03_AddProjectGitHubAccountFK(),
		schema.SCHEMA_20260822_04_AddProjectDokployServerFK(),
	})
	return migrator.Migrate()
}
