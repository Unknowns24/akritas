package migrations

import (
	"regexp"
	"sort"
	"testing"
)

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
		"20260822_10_add_investigations",
		"20260822_11_add_operations",
		"20260822_12_add_evidence",
		"20260823_01_add_monitoring_checkpoints",
		"20260823_02_add_incidents",
		"20260823_03_add_log_events",
		"20260823_04_link_investigation_history",
		"20260823_05_add_investigation_evidence_ids",
		"20260823_06_add_remediations",
		"20260823_07_add_validation_results",
		"20260823_06_add_github_issue_references",
		"20260823_07_enforce_issue_reference_investigation_incident",
		"20260823_08_add_dokploy_compose_sources",
		"20260823_09_allow_truthful_validation_output_redacted",
		"20260823_10_extend_remediation_lifecycle",
		"20260823_11_add_runtime_settings",
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

func TestMigrationRegistryDocumentsHistoricalSameSlotConflicts(t *testing.T) {
	t.Parallel()

	allowed := map[string]map[string]struct{}{
		"20260823_06": {
			"20260823_06_add_remediations":            {},
			"20260823_06_add_github_issue_references": {},
		},
		"20260823_07": {
			"20260823_07_add_validation_results":                         {},
			"20260823_07_enforce_issue_reference_investigation_incident": {},
		},
	}

	bySlot := map[string][]string{}
	for _, migration := range All() {
		slot := migrationSlot(t, migration.ID)
		bySlot[slot] = append(bySlot[slot], migration.ID)
	}
	for slot, ids := range bySlot {
		if len(ids) <= 1 {
			continue
		}
		sort.Strings(ids)
		allowedIDs := allowed[slot]
		if len(allowedIDs) != len(ids) {
			t.Fatalf("migration slot %s has undocumented duplicate IDs: %v", slot, ids)
		}
		for _, id := range ids {
			if _, ok := allowedIDs[id]; !ok {
				t.Fatalf("migration slot %s has undocumented duplicate ID %s", slot, id)
			}
		}
	}
}

func TestMigrationRegistryPreservesDependencyOrder(t *testing.T) {
	t.Parallel()

	positions := map[string]int{}
	for index, migration := range All() {
		if _, duplicate := positions[migration.ID]; duplicate {
			t.Fatalf("duplicate migration ID: %s", migration.ID)
		}
		positions[migration.ID] = index
	}

	assertBefore := func(first, second string) {
		t.Helper()
		if positions[first] >= positions[second] {
			t.Fatalf("migration %s must run before %s", first, second)
		}
	}

	assertBefore("20260823_02_add_incidents", "20260823_06_add_github_issue_references")
	assertBefore("20260822_10_add_investigations", "20260823_06_add_github_issue_references")
	assertBefore("20260823_06_add_github_issue_references", "20260823_07_enforce_issue_reference_investigation_incident")
	assertBefore("20260823_06_add_remediations", "20260823_07_add_validation_results")
	assertBefore("20260823_07_add_validation_results", "20260823_10_extend_remediation_lifecycle")
}

func migrationSlot(t *testing.T, id string) string {
	t.Helper()
	pattern := regexp.MustCompile(`^(\d{8}_\d{2})_[a-z0-9_]+$`)
	match := pattern.FindStringSubmatch(id)
	if match == nil {
		t.Fatalf("invalid migration ID: %s", id)
	}
	return match[1]
}
