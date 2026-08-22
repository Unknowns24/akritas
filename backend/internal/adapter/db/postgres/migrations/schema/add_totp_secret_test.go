package schema_test

import (
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations/schema"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// On a brand-new database, migration 20260822_01's AutoMigrate already
// creates administrators.encrypted_totp_secret because it reflects the
// current (already-extended) model.Administrator struct -- this migration
// only does real work on a database whose administrators table predates
// the field. Re-running it when the column already exists must be a no-op,
// not an error.
func TestMigration03IsIdempotentWhenColumnAlreadyExists(t *testing.T) {
	db := dbtest.Connect(t) // migrations.Run already applied 01-04

	migration := schema.Migration20260822_03_AddTotpSecretToAdministrators()
	if err := migration.Migrate(db); err != nil {
		t.Fatalf("re-running migration 03 on a database that already has the column must be a no-op, got: %v", err)
	}
	if !db.Migrator().HasColumn(&model.Administrator{}, "EncryptedTOTPSecret") {
		t.Fatal("column must still exist after the idempotent re-run")
	}
}
