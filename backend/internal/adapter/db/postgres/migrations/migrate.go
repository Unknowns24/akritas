package migrations

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations/schema"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func All() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		schema.SCHEMA_20260822_01_AddGitHubAccounts(),
		schema.SCHEMA_20260822_02_AddDokployServers(),
		schema.SCHEMA_20260822_03_AddCredentials(),
		schema.SCHEMA_20260822_04_AddGitHubAppRegistrations(),
		schema.SCHEMA_20260822_05_AddGitHubAppBindings(),
		schema.SCHEMA_20260822_06_AddAdministrators(),
		schema.SCHEMA_20260822_07_AddPendingEnrollments(),
		schema.SCHEMA_20260822_08_AddAdministratorSessions(),
		schema.SCHEMA_20260822_09_AddProjects(),
		schema.SCHEMA_20260822_10_AddInvestigations(),
		schema.SCHEMA_20260822_11_AddOperations(),
	}
}

func Run(db *gorm.DB) error {
	return gormigrate.New(db, gormigrate.DefaultOptions, All()).Migrate()
}
