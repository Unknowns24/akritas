package migrations

import "testing"

func TestMigrationRegistryIsOrderedAndReversible(t *testing.T) {
	expected := []string{
		"20260822_01_add_github_accounts",
		"20260822_02_add_dokploy_servers",
		"20260822_03_add_credentials",
		"20260822_04_add_github_app_registrations",
		"20260822_05_add_github_app_bindings",
		"20260822_06_add_administrators",
		"20260822_07_add_pending_enrollments",
		"20260822_08_add_administrator_sessions",
		"20260822_09_add_projects",
	}
	migrations := All()
	if len(migrations) != len(expected) {
		t.Fatalf("migration count = %d, want %d", len(migrations), len(expected))
	}
	for index, migration := range migrations {
		if migration.ID != expected[index] || migration.Migrate == nil || migration.Rollback == nil {
			t.Fatalf("invalid migration at %d: %#v", index, migration)
		}
	}
}
