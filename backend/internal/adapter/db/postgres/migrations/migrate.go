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
		schema.SCHEMA_20260822_12_AddEvidence(),
		schema.SCHEMA_20260823_01_AddMonitoringCheckpoints(),
		schema.SCHEMA_20260823_02_AddIncidents(),
		schema.SCHEMA_20260823_03_AddLogEvents(),
		schema.SCHEMA_20260823_04_LinkInvestigationHistory(),
		schema.SCHEMA_20260823_05_AddInvestigationEvidenceIDs(),
		schema.SCHEMA_20260823_06_AddDokployComposeSources(),
	}
}

func Run(db *gorm.DB) error {
	return gormigrate.New(db, gormigrate.DefaultOptions, All()).Migrate()
}
