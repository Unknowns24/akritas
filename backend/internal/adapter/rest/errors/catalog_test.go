package errors

import "testing"

func TestRESTErrorCatalogOwnsRESTCodes(t *testing.T) {
	t.Parallel()

	for name, value := range Catalog() {
		if value == nil || len(value.Code) != 9 || value.Code[2] != '1' {
			t.Fatalf("%s has invalid REST code: %+v", name, value)
		}
	}
}
