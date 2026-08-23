package dokploy

import (
	"context"
	"encoding/json"
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
	records, err := client.FetchLogs(context.Background(), portsout.LogFetchRequest{Server: serverEntity, Application: application, Tail: 10000, Since: "all"})
	if err != nil || len(records) != 1 || records[0].Message != "ERROR failed" {
		t.Fatalf("FetchLogs() = %+v, %v", records, err)
	}
}
