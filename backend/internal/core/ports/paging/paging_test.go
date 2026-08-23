package paging

import (
	"testing"

	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestParamsIsTheUkerPaginationContract(t *testing.T) {
	t.Parallel()

	value := Params{Limit: 25, Sort: []ukerpagination.SortExpression{{Field: "id", Direction: ukerpagination.DirectionDesc}}}
	var ukerValue ukerpagination.Params = value
	if ukerValue.Limit != 25 || len(ukerValue.Sort) != 1 {
		t.Fatalf("unexpected aliased params: %+v", ukerValue)
	}
}

func TestSliceCarriesProviderBoundariesWithoutEncodingThem(t *testing.T) {
	t.Parallel()

	page := Slice[string]{Items: []string{"one"}, Total: 1, NextBoundary: map[string]string{"provider_page": "2"}}
	if page.NextBoundary["provider_page"] != "2" {
		t.Fatalf("unexpected boundary: %+v", page.NextBoundary)
	}
}
