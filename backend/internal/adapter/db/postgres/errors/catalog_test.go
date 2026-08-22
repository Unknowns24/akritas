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
