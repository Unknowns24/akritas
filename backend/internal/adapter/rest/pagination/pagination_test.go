package pagination

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestParseUsesUkerAndBindsCursorToIntegrationScope(t *testing.T) {
	t.Parallel()
	config, err := NewConfig([]byte("01234567890123456789012345678901"), 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	policy := Policy{AllowedFilters: map[string]struct{}{"name_like": {}}, Scope: "github:account-1", BoundaryField: "provider_page"}
	params, err := Parse(httptest.NewRequest("GET", "/?limit=25&name_like=api", nil), config, policy)
	if err != nil {
		t.Fatal(err)
	}
	page, err := BuildProviderPage(params, providerResult(), config)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("GET", "/?cursor="+page.Paging.NextCursor, nil)
	if _, err := Parse(request, config, Policy{AllowedFilters: policy.AllowedFilters, Scope: "github:account-2", BoundaryField: "provider_page"}); err == nil {
		t.Fatal("cursor for another integration scope must be rejected")
	}
}

func TestParseRejectsTamperingExpiryAndCursorOverrides(t *testing.T) {
	t.Parallel()
	secret := []byte("01234567890123456789012345678901")
	config, _ := NewConfig(secret, 15*time.Minute)
	policy := Policy{AllowedFilters: map[string]struct{}{"name_like": {}}, Scope: "github:account-1", BoundaryField: "provider_page"}
	params, err := Parse(httptest.NewRequest("GET", "/?name_like=api", nil), config, policy)
	if err != nil {
		t.Fatal(err)
	}
	page, err := BuildProviderPage(params, providerResult(), config)
	if err != nil {
		t.Fatal(err)
	}
	cursor := page.Paging.NextCursor
	last := byte('A')
	if cursor[len(cursor)-1] == last {
		last = 'B'
	}
	if _, err := Parse(httptest.NewRequest("GET", "/?cursor="+cursor[:len(cursor)-1]+string(last), nil), config, policy); err == nil {
		t.Fatal("tampered cursor must be rejected")
	}
	if _, err := Parse(httptest.NewRequest("GET", "/?cursor="+cursor+"&name_like=other", nil), config, policy); err == nil {
		t.Fatal("cursor filters must not be overridden")
	}

	expired, err := ukerpagination.EncodeCursorSigned(ukerpagination.CursorPayload{
		Limit:     25,
		Sort:      []ukerpagination.SortExpression{{Field: "provider_page", Direction: ukerpagination.DirectionAsc}},
		Filters:   map[string]string{"integration_scope_eq": "github:account-1"},
		After:     map[string]string{"provider_page": "2"},
		Timestamp: time.Now().Add(-time.Hour).Unix(),
	}, secret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Parse(httptest.NewRequest("GET", "/?cursor="+expired, nil), config, policy)
	if !errors.Is(err, ukerpagination.ErrCursorExpired) {
		t.Fatalf("expired cursor error = %v", err)
	}
}

func TestParseRejectsUnknownFilterAndSort(t *testing.T) {
	t.Parallel()
	config, _ := NewConfig([]byte("01234567890123456789012345678901"), time.Minute)
	policy := Policy{AllowedFilters: map[string]struct{}{"name_like": {}}, AllowedSorts: map[string]struct{}{"created_at": {}, "id": {}}}
	for _, target := range []string{"/?owner_id_eq=secret", "/?sort=updated_at:desc"} {
		if _, err := Parse(httptest.NewRequest("GET", target, nil), config, policy); err == nil {
			t.Fatalf("disallowed query %q accepted", target)
		}
	}
}

func TestParseAppendsStableIDTieBreaker(t *testing.T) {
	t.Parallel()
	config, _ := NewConfig([]byte("01234567890123456789012345678901"), time.Minute)
	params, err := Parse(httptest.NewRequest("GET", "/?sort=created_at:desc", nil), config, Policy{AllowedSorts: map[string]struct{}{"created_at": {}, "id": {}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(params.Sort) != 2 || params.Sort[1].Field != "id" || params.Sort[1].Direction != ukerpagination.DirectionDesc {
		t.Fatalf("unstable sort: %+v", params.Sort)
	}
}

func providerResult() paging.Slice[string] {
	return paging.Slice[string]{Items: []string{"api"}, Total: 2, NextBoundary: map[string]string{"provider_page": "2"}}
}
