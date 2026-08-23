// Package dbtest provides a real Postgres connection for repository tests,
// per the approved AKR-AUTH-BOOTSTRAP TDD plan (repository tests run
// against local Postgres, not sqlmock).
package dbtest

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
)

const defaultDSN = "postgres://localhost:5432/akritas_test?sslmode=disable"

// Connect opens a migrated connection to the local test database and
// truncates every table it manages so each test starts from a clean slate.
// It skips the test if no database is reachable.
func Connect(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("AKRITAS_TEST_DB_DSN")
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := postgres.Open(postgres.Config{DSN: dsn})
	if err != nil {
		t.Skipf("skipping: local Postgres not reachable at %q: %v", dsn, err)
	}

	if err := migrations.Run(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE projects, administrator_sessions, pending_enrollments, administrators, github_app_bindings, github_app_registrations, credentials, github_accounts, dokploy_servers CASCADE").Error; err != nil {
		t.Fatalf("truncate tables: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}
