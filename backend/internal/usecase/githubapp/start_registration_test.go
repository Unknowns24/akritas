package githubapp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func TestStartRegistrationFormActionIncludesState(t *testing.T) {
	const state = "conversion-state+value/with?reserved&chars=0001"
	tests := []struct {
		name                string
		command             portsin.StartGitHubAppRegistrationCommand
		expectedEscapedPath string
	}{
		{
			name: "personal account",
			command: portsin.StartGitHubAppRegistrationCommand{
				DisplayName: "Personal",
				OwnerType:   domain.GitHubAccountPersonal,
			},
			expectedEscapedPath: "/settings/apps/new",
		},
		{
			name: "organization account",
			command: portsin.StartGitHubAppRegistrationCommand{
				DisplayName:  "Organization",
				OwnerType:    domain.GitHubAccountOrganization,
				Organization: "Unknowns24",
			},
			expectedEscapedPath: "/organizations/Unknowns24/settings/apps/new",
		},
		{
			name: "escaped organization path",
			command: portsin.StartGitHubAppRegistrationCommand{
				DisplayName:  "Escaped organization",
				OwnerType:    domain.GitHubAccountOrganization,
				Organization: "Unknowns24/team name",
			},
			expectedEscapedPath: "/organizations/Unknowns24%2Fteam%20name/settings/apps/new",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &registrationStoreFake{}
			uc, err := New(store, &appGatewayFake{}, "https://akritas.example.com", uuid.New, func() time.Time {
				return time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
			})
			if err != nil {
				t.Fatal(err)
			}
			uc.newState = func() (string, error) { return state, nil }

			started, err := uc.StartRegistration(context.Background(), tt.command)
			if err != nil {
				t.Fatal(err)
			}
			formURL, err := url.Parse(started.FormAction)
			if err != nil {
				t.Fatalf("parse form action: %v", err)
			}
			if formURL.Scheme != "https" || formURL.Host != "github.com" || formURL.EscapedPath() != tt.expectedEscapedPath {
				t.Fatalf("unexpected form action endpoint: %q", started.FormAction)
			}
			query := formURL.Query()
			if len(query) != 1 || len(query["state"]) != 1 || query.Get("state") != state {
				t.Fatalf("form action must contain exactly the generated state: %q", started.FormAction)
			}
			if formURL.RawQuery != "state="+url.QueryEscape(state) {
				t.Fatalf("state is not encoded canonically: %q", formURL.RawQuery)
			}
			if started.State != state {
				t.Fatalf("compatibility state changed: got %q", started.State)
			}
			expectedDigest := sha256.Sum256([]byte(state))
			if len(store.registration.ConversionStateDigest) != sha256.Size || !bytes.Equal(store.registration.ConversionStateDigest, expectedDigest[:]) {
				t.Fatal("registration did not persist exactly sha256(state)")
			}
			if bytes.Equal(store.registration.ConversionStateDigest, []byte(state)) {
				t.Fatal("raw state was persisted")
			}
		})
	}
}

func TestNewValidatesPublicURL(t *testing.T) {
	tests := []struct {
		name      string
		publicURL string
		wantErr   bool
	}{
		{name: "public HTTPS", publicURL: "https://akritas.example.com"},
		{name: "localhost HTTP", publicURL: "http://localhost:8080"},
		{name: "IPv4 loopback HTTP", publicURL: "http://127.0.0.1:8080"},
		{name: "IPv6 loopback HTTP", publicURL: "http://[::1]:8080"},
		{name: "public HTTP", publicURL: "http://akritas.example.com", wantErr: true},
		{name: "wildcard HTTP", publicURL: "http://0.0.0.0:8080", wantErr: true},
		{name: "public URL with path", publicURL: "https://akritas.example.com/base", wantErr: true},
		{name: "public URL with query", publicURL: "https://akritas.example.com?next=/settings", wantErr: true},
		{name: "public URL with credentials", publicURL: "https://user@akritas.example.com", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(&registrationStoreFake{}, &appGatewayFake{}, tt.publicURL, uuid.New, time.Now)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidConfiguration) {
					t.Fatalf("expected invalid configuration, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("valid public URL rejected: %v", err)
			}
		})
	}
}

func TestCompleteManifestRejectsMissingAndExpiredStateBeforeExchange(t *testing.T) {
	const state = "conversion-state-value-that-is-long-enough-0001"
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	currentTime := now
	store := &registrationStoreFake{}
	gateway := &appGatewayFake{}
	uc, err := New(store, gateway, "https://akritas.example.com", uuid.New, func() time.Time { return currentTime })
	if err != nil {
		t.Fatal(err)
	}
	uc.newState = func() (string, error) { return state, nil }

	started, err := uc.StartRegistration(context.Background(), portsin.StartGitHubAppRegistrationCommand{
		DisplayName: "Personal",
		OwnerType:   domain.GitHubAccountPersonal,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := uc.CompleteManifest(context.Background(), "manifest-code-at-least-twenty", ""); !errors.Is(err, domain.ErrManifestStateInvalid) {
		t.Fatalf("expected invalid missing state, got %v", err)
	}
	currentTime = now.Add(time.Hour + time.Second)
	if _, err := uc.CompleteManifest(context.Background(), "manifest-code-at-least-twenty", started.State); !errors.Is(err, domain.ErrManifestStateConflict) {
		t.Fatalf("expected expired state conflict, got %v", err)
	}
	if gateway.exchangeCalls != 0 {
		t.Fatalf("invalid callbacks reached GitHub exchange %d times", gateway.exchangeCalls)
	}
}
