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
	} {
		if catalog[name] == nil {
			t.Fatalf("merged PostgreSQL catalog is missing %s", name)
		}
	}
}
