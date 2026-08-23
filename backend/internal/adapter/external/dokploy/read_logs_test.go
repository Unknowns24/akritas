package dokploy

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func TestParseLogsPreservesTimestampOrdinalAndMultilinePayload(t *testing.T) {
	body := `"2026-08-22T12:00:00.123456789Z first\ncontinuation\n2026-08-22T12:00:00.123456789Z second\n2026-08-22T12:00:01Z third\n"`
	got, err := ParseLogs([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("records = %d, want 3: %+v", len(got), got)
	}
	if got[0].Message != "first\ncontinuation" || got[0].Ordinal != 0 || got[1].Ordinal != 1 || got[2].Ordinal != 0 {
		t.Fatalf("unexpected records: %+v", got)
	}
	if !got[0].Timestamp.Equal(time.Date(2026, 8, 22, 12, 0, 0, 123456789, time.UTC)) {
		t.Fatalf("timestamp = %s", got[0].Timestamp)
	}
	if got[0].ContentHash == "" || got[0].ContentHash == got[1].ContentHash {
		t.Fatal("content hashes are not stable/distinct")
	}
}

func TestFetchLogsUsesExistingCredentialTransportAndPublishedReadLogsQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/application.readLogs" || r.URL.Query().Get("applicationId") != "app-1" || r.URL.Query().Get("tail") != "10000" || r.URL.Query().Get("since") != "all" {
			t.Fatalf("request = %s", r.URL.String())
		}
		if r.Header.Get("x-api-key") != "dokploy-key" {
			t.Fatal("missing shared credential header")
		}
		_ = json.NewEncoder(w).Encode("2026-08-22T12:00:00Z ERROR failed\n")
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	application, err := domain.NewDokployApplication(serverEntity.ID, "app-1", "instance-1", "API", "production", domain.DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("dokploy-key")})
	source, err := domain.SourceFromApplication(application)
	if err != nil {
		t.Fatal(err)
	}
	records, err := client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Source: source, Tail: 10000, Since: "all"})
	if err != nil || len(records) != 1 || records[0].Message != "ERROR failed" {
		t.Fatalf("FetchLogs() = %+v, %v", records, err)
	}
}

func TestFetchComposeLogsResolvesDeterministicRunningContainer(t *testing.T) {
	requests := make([]string, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.String())
		switch r.URL.Path {
		case "/api/docker.getContainersByAppNameMatch":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"Id": "bbb", "State": "running", "Labels": map[string]string{"com.docker.compose.project": "stack-prod", "com.docker.compose.service": "api"}},
				{"Id": "aaa", "State": "running", "Labels": map[string]string{"com.docker.compose.project": "stack-prod", "com.docker.compose.service": "api"}},
				{"Id": "000", "State": "exited", "Labels": map[string]string{"com.docker.compose.project": "stack-prod", "com.docker.compose.service": "api"}},
			})
		case "/api/compose.readLogs":
			if r.URL.Query().Get("containerId") != "aaa" || r.URL.Query().Get("composeId") != "compose-1" {
				t.Fatalf("read logs query = %s", r.URL.String())
			}
			_ = json.NewEncoder(w).Encode("2026-08-22T12:00:00Z ready\n")
		default:
			t.Fatalf("unexpected request %s", r.URL.String())
		}
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	source, err := domain.NewDokploySource(serverEntity.ID, domain.DokploySourceComposeService, "compose-1", "api", "stack-prod", "Platform / api", "production", domain.DokploySourceRunning, domain.DokployRuntimeCompose, "remote-1")
	if err != nil {
		t.Fatal(err)
	}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	records, err := client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Source: source, Tail: 100, Since: "1h"})
	if err != nil || len(records) != 1 || records[0].Message != "ready" || len(requests) != 2 {
		t.Fatalf("FetchLogs() = %+v, %v, requests=%v", records, err, requests)
	}
}

func TestFetchComposeLogsMatchesServiceWhenProjectLabelDiffersFromDokployAppName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/docker.getContainersByAppNameMatch":
			if r.URL.Query().Get("appName") != "lervu-frontend-2asbvo" || r.URL.Query().Get("appType") != "docker-compose" {
				t.Fatalf("container query = %s", r.URL.String())
			}
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"Id": "frontend-container", "State": "running", "Labels": map[string]string{"com.docker.compose.project": "lervu-frontend", "com.docker.compose.service": "lervu_frontend"}},
			})
		case "/api/compose.readLogs":
			if r.URL.Query().Get("containerId") != "frontend-container" {
				t.Fatalf("read logs query = %s", r.URL.String())
			}
			_ = json.NewEncoder(w).Encode("2026-08-22T12:00:00Z ready\n")
		default:
			t.Fatalf("unexpected request %s", r.URL.String())
		}
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	source, err := domain.NewDokploySource(serverEntity.ID, domain.DokploySourceComposeService, "compose-1", "lervu_frontend", "lervu-frontend-2asbvo", "frontend / lervu_frontend", "production", domain.DokploySourceUnknown, domain.DokployRuntimeCompose, "")
	if err != nil {
		t.Fatal(err)
	}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	records, err := client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Source: source, Tail: 100, Since: "1h"})
	if err != nil || len(records) != 1 || records[0].Message != "ready" {
		t.Fatalf("FetchLogs() = %+v, %v", records, err)
	}
}

func TestFetchComposeLogsFallsBackToSingleRunningContainerWhenServiceLabelDiffers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/docker.getContainersByAppNameMatch":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"id": "only-running-container", "state": "running", "labels": map[string]string{"com.docker.compose.service": "backend"}},
			})
		case "/api/compose.readLogs":
			if r.URL.Query().Get("containerId") != "only-running-container" {
				t.Fatalf("read logs query = %s", r.URL.String())
			}
			_ = json.NewEncoder(w).Encode("2026-08-22T12:00:00Z ready\n")
		default:
			t.Fatalf("unexpected request %s", r.URL.String())
		}
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	source, err := domain.NewDokploySource(serverEntity.ID, domain.DokploySourceComposeService, "compose-1", "akritas-backend-error", "hackathonprueba-backendwitherrors-esmjqu", "backend / akritas-backend-error", "", domain.DokploySourceUnknown, domain.DokployRuntimeCompose, "")
	if err != nil {
		t.Fatal(err)
	}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	records, err := client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Source: source, Tail: 100, Since: "1h"})
	if err != nil || len(records) != 1 || records[0].Message != "ready" {
		t.Fatalf("FetchLogs() = %+v, %v", records, err)
	}
}

func TestFetchComposeLogsDoesNotFallbackWhenMultipleRunningContainersDiffer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/docker.getContainersByAppNameMatch":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"Id": "api-container", "State": "running", "Labels": map[string]string{"com.docker.compose.service": "api"}},
				{"Id": "worker-container", "State": "running", "Labels": map[string]string{"com.docker.compose.service": "worker"}},
			})
		default:
			t.Fatalf("unexpected request %s", r.URL.String())
		}
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	source, err := domain.NewDokploySource(serverEntity.ID, domain.DokploySourceComposeService, "compose-1", "missing", "stack-prod", "Platform / missing", "", domain.DokploySourceUnknown, domain.DokployRuntimeCompose, "")
	if err != nil {
		t.Fatal(err)
	}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	_, err = client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Source: source, Tail: 100, Since: "1h"})
	if !errors.Is(err, domain.ErrDokployContainerUnavailable) {
		t.Fatalf("FetchLogs() error = %v", err)
	}
}

func TestFetchComposeLogsReadsSingleStoppedContainer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/docker.getContainersByAppNameMatch":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"id": "crashed-container", "state": "exited"},
			})
		case "/api/compose.readLogs":
			if r.URL.Query().Get("containerId") != "crashed-container" {
				t.Fatalf("read logs query = %s", r.URL.String())
			}
			_ = json.NewEncoder(w).Encode("2026-08-22T12:00:00Z ERROR crashed\n")
		default:
			t.Fatalf("unexpected request %s", r.URL.String())
		}
	}))
	defer server.Close()
	serverEntity := dokployServer(t, server.URL)
	source, err := domain.NewDokploySource(serverEntity.ID, domain.DokploySourceComposeService, "compose-1", "akritas-backend-error", "hackathonprueba-backendwitherrors-esmjqu", "backend / akritas-backend-error", "", domain.DokploySourceUnknown, domain.DokployRuntimeCompose, "")
	if err != nil {
		t.Fatal(err)
	}
	client := newDokployTestClientWithCredentials(t, server.URL, resolverFake{ips: []net.IP{net.ParseIP("127.0.0.1")}}, dokployCredentialStoreFake{value: []byte("key")})
	records, err := client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Source: source, Tail: 100, Since: "1h"})
	if err != nil || len(records) != 1 || records[0].Message != "ERROR crashed" {
		t.Fatalf("FetchLogs() = %+v, %v", records, err)
	}
}
