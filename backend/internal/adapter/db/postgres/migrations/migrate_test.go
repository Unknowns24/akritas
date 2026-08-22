package migrations_test

import (
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// dbtest.Connect already runs migrations.Run; this test only asserts the
// resulting schema is queryable for every table/column it manages.
func TestMigrationsCreateExpectedTables(t *testing.T) {
	db := dbtest.Connect(t)

	if !db.Migrator().HasTable(&model.Administrator{}) {
		t.Fatal("administrators table was not created")
	}
	if !db.Migrator().HasTable(&model.PendingEnrollment{}) {
		t.Fatal("pending_enrollments table was not created")
	}
	if !db.Migrator().HasColumn(&model.Administrator{}, "EncryptedTOTPSecret") {
		t.Fatal("administrators.encrypted_totp_secret column was not created")
	}
	if !db.Migrator().HasTable(&model.AdministratorSession{}) {
		t.Fatal("administrator_sessions table was not created")
	}
}
