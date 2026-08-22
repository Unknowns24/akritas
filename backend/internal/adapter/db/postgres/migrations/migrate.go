package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations/schema"
)

// Run applies every schema migration, in order, using gormigrate. AutoMigrate
// is never called globally — only inside a specific versioned migration.
func Run(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		schema.Migration20260822_01_CreateAdministrators(),
		schema.Migration20260822_02_CreatePendingEnrollments(),
		schema.Migration20260822_03_AddTotpSecretToAdministrators(),
		schema.Migration20260822_04_CreateAdministratorSessions(),
		schema.Migration20260822_05_AddLastAcceptedTotpPeriodToAdministrators(),
	})
	return m.Migrate()
}
