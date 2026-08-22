package schema_test

import (
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations/schema"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// Same lesson as migration 20260822_03: on a brand-new database, migration
// 20260822_01's AutoMigrate already creates this column because it reflects
// the current (already-extended) model.Administrator struct. Re-running
// this migration when the column already exists must be a no-op.
func TestMigration05IsIdempotentWhenColumnAlreadyExists(t *testing.T) {
	db := dbtest.Connect(t) // migrations.Run already applied 01-05

	migration := schema.Migration20260822_05_AddLastAcceptedTotpPeriodToAdministrators()
	if err := migration.Migrate(db); err != nil {
		t.Fatalf("re-running migration 05 on a database that already has the column must be a no-op, got: %v", err)
	}
	if !db.Migrator().HasColumn(&model.Administrator{}, "LastAcceptedTOTPPeriod") {
		t.Fatal("column must still exist after the idempotent re-run")
	}
}
