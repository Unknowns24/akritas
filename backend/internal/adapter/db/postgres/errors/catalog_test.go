package errors

import "testing"

func TestPostgresErrorCatalogOwnsDatabaseCodes(t *testing.T) {
	t.Parallel()

	for name, value := range Catalog() {
		if value == nil || len(value.Code) != 9 || value.Code[2] != '2' {
			t.Fatalf("%s has invalid DB code: %+v", name, value)
		}
	}
}

func TestPostgresErrorCatalogContainsMergedH1H2H3Sentinels(t *testing.T) {
	t.Parallel()
	catalog := Catalog()
	for _, name := range []string{
		"ErrIntegrationPersistence", "ErrProjectPersistence", "ErrInvestigationPersistence",
		"ErrOperationPersistence", "ErrEvidencePersistence", "ErrIncidentPersistence", "ErrMonitoringPersistence",
		"ErrRemediationPersistence", "ErrValidationResultPersistence", "ErrGitHubIssueReferencePersistence",
	} {
		if catalog[name] == nil {
			t.Fatalf("merged PostgreSQL catalog is missing %s", name)
		}
	}
}

func TestPostgresErrorCatalogHasDistinctCodes(t *testing.T) {
	t.Parallel()

	seen := map[string]string{}
	for name, value := range Catalog() {
		if value == nil {
			t.Fatalf("%s is nil", name)
		}
		if previous := seen[value.Code]; previous != "" {
			t.Fatalf("duplicate code %s for %s and %s", value.Code, previous, name)
		}
		seen[value.Code] = name
	}
}
