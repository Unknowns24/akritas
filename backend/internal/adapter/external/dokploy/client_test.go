package dokploy

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestValidateNormalizesURLAndUsesAPIKeyHeader(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/settings.health" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "dokploy-key" {
			t.Fatal("missing x-api-key")
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()
	client := newDokployTestClient(t, server.URL, nil)

	result, err := client.Validate(context.Background(), server.URL+"/", "dokploy-key")
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if result.NormalizedBaseURL != server.URL || len(result.ServerIdentifier) != 64 {
		t.Fatalf("Validate() = %+v", result)
	}
}

func TestValidateRejectsUnsafeHTTPDestination(t *testing.T) {
	t.Parallel()

	client := newDokployTestClient(t, "http://127.0.0.1", resolverFake{ips: []net.IP{net.ParseIP("169.254.169.254")}})
	_, err := client.Validate(context.Background(), "http://metadata.internal", "dokploy-key")
	if err == nil {
		t.Fatal("Validate() accepted link-local metadata destination")
	}
}

func TestRequestDialRejectsDNSRebindingToMetadata(t *testing.T) {
	t.Parallel()
	resolver := &sequenceResolver{answers: [][]net.IP{
		{net.ParseIP("127.0.0.1")},
		{net.ParseIP("127.0.0.1")},
		{net.ParseIP("169.254.169.254")},
	}}
	client := newDokployTestClient(t, "http://127.0.0.1", resolver)
	_, err := client.Validate(context.Background(), "http://dokploy.internal:3000", "dokploy-key")
	if err == nil {
		t.Fatal("DNS rebinding to metadata was accepted")
	}
}

func TestProviderResponseSizeIsBounded(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", maximumBodyBytes+1)))
	}))
	defer server.Close()
	client := newDokployTestClient(t, server.URL, nil)
	if _, err := client.Validate(context.Background(), server.URL, "dokploy-key"); err == nil {
		t.Fatal("oversized Dokploy response was accepted")
	}
}

func TestListApplicationsMapsProviderPayloadAndUnknownStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/application.search" || r.URL.Query().Get("limit") != "20" || r.URL.Query().Get("offset") != "40" {
			t.Fatalf("request = %s", r.URL.String())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{"applicationId": "app-1", "appName": "api", "name": "API", "status": "brand_new", "environment": map[string]any{"name": "production"}}},
			"total": 1,
		})
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	credentials := dokployCredentialStoreFake{value: []byte("dokploy-key")}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, credentials)

	page, err := client.ListApplications(context.Background(), serverEntity, paging.Params{
		Limit:  20,
		Cursor: &ukerpagination.CursorPayload{After: map[string]string{"provider_offset": "40"}},
	})
	if err != nil {
		t.Fatalf("ListApplications() error = %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].ApplicationIdentifier != "app-1" || page.Items[0].InstanceIdentifier != "api" || page.Items[0].Environment != "production" || page.Items[0].Status != domain.DokployApplicationUnknown {
		t.Fatalf("ListApplications() = %+v", page)
	}
	if page.PrevBoundary["provider_offset"] != "20" {
		t.Fatalf("provider boundary was not translated: %+v", page.PrevBoundary)
	}
}

func TestGetApplicationRequiresExactIdentifier(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/application.search" || r.URL.Query().Get("q") != "app-1" {
			t.Fatalf("request = %s", r.URL.String())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"items": []map[string]any{{"applicationId": "app-10", "appName": "wrong", "name": "Wrong"}, {"applicationId": "app-1", "appName": "api", "name": "API", "status": "running", "environment": "production"}}})
	}))
	defer server.Close()
	entity := dokployServer(t, server.URL)
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	application, err := client.GetApplication(context.Background(), entity, "app-1")
	if err != nil || application.ApplicationIdentifier != "app-1" || application.DisplayName != "API" || application.Status != domain.DokployApplicationRunning {
		t.Fatalf("GetApplication() = %+v, %v", application, err)
	}
}

func TestGetApplicationNormalizesCredentialRejection(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "api key must not leak", http.StatusUnauthorized)
	}))
	defer server.Close()
	entity := dokployServer(t, server.URL)
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	if _, err := client.GetApplication(context.Background(), entity, "app-1"); !errors.Is(err, domain.ErrDokployCredentialRejected) || strings.Contains(err.Error(), "api key") {
		t.Fatalf("credential rejection = %v", err)
	}
}

func newDokployTestClient(t *testing.T, rawURL string, resolver Resolver) *Client {
	return newDokployTestClientWithCredentials(t, rawURL, resolver, dokployCredentialStoreFake{})
}

func newDokployTestClientWithCredentials(t *testing.T, _ string, resolver Resolver, credentials portsout.CredentialStore) *Client {
	t.Helper()
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	client, err := NewClient(ClientConfig{HTTPClient: &http.Client{Timeout: time.Second}, Resolver: resolver, Credentials: credentials})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func dokployServer(t *testing.T, rawURL string) domain.DokployServer {
	t.Helper()
	server, err := domain.NewDokployServer(uuid.New(), "Dokploy", strings.TrimSuffix(rawURL, "/"), strings.Repeat("a", 64), domain.IntegrationStatusConnected, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	server.CredentialConfigured = true
	return *server
}

type resolverFake struct{ ips []net.IP }

func (r resolverFake) LookupIP(_ context.Context, _, _ string) ([]net.IP, error) { return r.ips, nil }

type sequenceResolver struct {
	mu      sync.Mutex
	answers [][]net.IP
}

func (r *sequenceResolver) LookupIP(_ context.Context, _, _ string) ([]net.IP, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.answers) == 0 {
		return []net.IP{net.ParseIP("169.254.169.254")}, nil
	}
	answer := r.answers[0]
	r.answers = r.answers[1:]
	return answer, nil
}

type dokployCredentialStoreFake struct{ value []byte }

func (f dokployCredentialStoreFake) Put(context.Context, string, uuid.UUID, portsout.SecretValue) error {
	return nil
}
func (f dokployCredentialStoreFake) Get(context.Context, string, uuid.UUID, portsout.SecretKind) ([]byte, error) {
	return append([]byte(nil), f.value...), nil
}
func (f dokployCredentialStoreFake) DeleteOwner(context.Context, string, uuid.UUID) error { return nil }
func (f dokployCredentialStoreFake) MoveOwner(context.Context, string, uuid.UUID, string, uuid.UUID) error {
	return nil
}
