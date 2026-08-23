package github

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestValidatePATUsesVersionedAPIAndValidatesPersonalIdentity(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer github_pat_test_value" {
			t.Fatal("missing bearer token")
		}
		if r.Header.Get("X-GitHub-Api-Version") != APIVersion || r.Header.Get("Accept") == "" {
			t.Fatal("missing GitHub versioned headers")
		}
		w.Header().Set("X-OAuth-Scopes", "repo, read:org")
		_ = json.NewEncoder(w).Encode(map[string]any{"login": "Unknowns24", "type": "User"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL, nil)
	result, err := client.ValidatePAT(context.Background(), portsout.GitHubPATValidationRequest{
		AccountType: domain.GitHubAccountPersonal, AccountIdentifier: "unknowns24", Token: "github_pat_test_value",
	})
	if err != nil {
		t.Fatalf("ValidatePAT() error = %v", err)
	}
	if result.AccountIdentifier != "Unknowns24" || len(result.ClassicScopes) != 2 {
		t.Fatalf("ValidatePAT() = %+v", result)
	}
}

func TestGitHubAppVerifiesInstallationAndUsesEphemeralInstallationToken(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	accountID := uuid.New()
	registrationID := uuid.New()
	accessTokenCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/app/installations/99":
			if parts := len(splitToken(r.Header.Get("Authorization"))); parts != 3 {
				t.Fatalf("expected App JWT, got %q", r.Header.Get("Authorization"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 99, "account": map[string]any{"login": "acme", "type": "Organization"}})
		case "/app/installations/99/access_tokens":
			accessTokenCalls++
			_ = json.NewEncoder(w).Encode(map[string]any{"token": "ephemeral-installation-token", "expires_at": now.Add(time.Hour)})
		case "/installation/repositories":
			if r.Header.Get("Authorization") != "Bearer ephemeral-installation-token" {
				t.Fatalf("installation token not used: %q", r.Header.Get("Authorization"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"total_count": 1, "repositories": []map[string]any{{
				"id": 1, "name": "service", "default_branch": "main", "private": true, "html_url": "https://github.com/acme/service", "owner": map[string]any{"login": "acme"},
			}}})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()
	base, _ := url.Parse(server.URL)
	credentials := credentialStoreFake{values: map[string][]byte{
		credentialKey(registrationID, portsout.SecretKindGitHubPrivateKey): privateKeyPEM,
		credentialKey(accountID, portsout.SecretKindGitHubPrivateKey):      privateKeyPEM,
	}}
	bindings := bindingReaderFake{binding: portsout.GitHubAppBinding{GitHubAccountID: accountID, AppID: 42, InstallationID: 99, AppSlug: "akritas"}}
	client, err := NewClient(ClientConfig{APIBaseURL: base, HTTPClient: &http.Client{Timeout: time.Second}, Credentials: credentials, Bindings: bindings, Now: func() time.Time { return now }})
	if err != nil {
		t.Fatal(err)
	}
	appID := int64(42)
	installation, err := client.VerifyInstallation(context.Background(), portsout.GitHubAppRegistration{ID: registrationID, AppID: &appID}, 99)
	if err != nil || installation.AccountLogin != "acme" || installation.AccountType != domain.GitHubAccountOrganization {
		t.Fatalf("installation verification failed: %+v %v", installation, err)
	}
	account, err := domain.NewGitHubAccount(accountID, "Acme", domain.GitHubAccountOrganization, domain.GitHubAuthenticationGitHubApp, "acme", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := client.ListRepositories(context.Background(), *account, paging.Params{Limit: 20}); err != nil {
			t.Fatal(err)
		}
	}
	if accessTokenCalls != 1 {
		t.Fatalf("installation token was not cached only in memory: calls=%d", accessTokenCalls)
	}
}

func splitToken(value string) []string {
	const prefix = "Bearer "
	if len(value) < len(prefix) || value[:len(prefix)] != prefix {
		return nil
	}
	result := make([]string, 0, 3)
	start := len(prefix)
	for index := start; index <= len(value); index++ {
		if index == len(value) || value[index] == '.' {
			result = append(result, value[start:index])
			start = index + 1
		}
	}
	return result
}

func TestConnectionTestNormalizesProviderAuthenticationFailure(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `token github_pat_should_not_leak`, http.StatusUnauthorized)
	}))
	defer server.Close()
	account := githubAccount(t)
	credentials := credentialStoreFake{values: map[string][]byte{credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("secret")}}
	client := newTestClient(t, server.URL, credentials)

	result, err := client.TestConnection(context.Background(), account)
	if err != nil {
		t.Fatalf("TestConnection() error = %v", err)
	}
	if result.Status != domain.ConnectionTestAuthenticationFailed || result.CheckedAt.IsZero() || result.Latency < 0 {
		t.Fatalf("TestConnection() = %+v", result)
	}
	if result.UserMessage == "" || result.UserMessage == "token github_pat_should_not_leak" {
		t.Fatal("TestConnection() leaked provider response")
	}
}

func TestProviderResponseSizeIsBoundedWithoutLeakingBody(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", maximumBodyBytes+1)))
	}))
	defer server.Close()
	client := newTestClient(t, server.URL, nil)
	_, err := client.ValidatePAT(context.Background(), portsout.GitHubPATValidationRequest{AccountType: domain.GitHubAccountPersonal, AccountIdentifier: "acme", Token: "github_pat_test_value"})
	if err == nil || strings.Contains(err.Error(), "xxxxx") {
		t.Fatalf("oversized provider body was accepted or leaked: %v", err)
	}
}

func TestListRepositoriesMapsOnlyConfiguredOwner(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user/repos" || r.URL.Query().Get("per_page") != "20" || r.URL.Query().Get("page") != "3" {
			t.Fatalf("request = %s", r.URL.String())
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "akritas", "full_name": "Unknowns24/akritas", "default_branch": "main", "private": true, "html_url": "https://github.com/Unknowns24/akritas", "owner": map[string]any{"login": "Unknowns24"}},
			{"id": 2, "name": "other", "full_name": "another/other", "default_branch": "main", "private": false, "html_url": "https://github.com/another/other", "owner": map[string]any{"login": "another"}},
		})
	}))
	defer server.Close()
	account := githubAccount(t)
	credentials := credentialStoreFake{values: map[string][]byte{credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("secret")}}
	client := newTestClient(t, server.URL, credentials)

	page, err := client.ListRepositories(context.Background(), account, paging.Params{
		Limit:   20,
		Filters: map[string]string{"name_like": "ak"},
		Cursor:  &ukerpagination.CursorPayload{After: map[string]string{"provider_page": "3"}},
	})
	if err != nil {
		t.Fatalf("ListRepositories() error = %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].FullName != "Unknowns24/akritas" || page.Items[0].RepositoryIdentifier != "1" {
		t.Fatalf("ListRepositories() = %+v", page)
	}
	if page.PrevBoundary["provider_page"] != "2" {
		t.Fatalf("provider boundary was not translated: %+v", page.PrevBoundary)
	}
}

func TestGetRepositoryResolvesOpaqueIDAndRejectsAnotherOwner(t *testing.T) {
	t.Parallel()
	owner := "Unknowns24"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repositories/42" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 42, "name": "akritas", "default_branch": "main", "private": true, "html_url": "https://github.com/" + owner + "/akritas", "owner": map[string]any{"login": owner}})
	}))
	defer server.Close()
	account := githubAccount(t)
	credentials := credentialStoreFake{values: map[string][]byte{credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("secret")}}
	client := newTestClient(t, server.URL, credentials)
	repository, err := client.GetRepository(context.Background(), account, "42")
	if err != nil || repository.RepositoryIdentifier != "42" || repository.DefaultBranch != "main" || !repository.Private {
		t.Fatalf("GetRepository() = %+v, %v", repository, err)
	}
	owner = "another"
	if _, err := client.GetRepository(context.Background(), account, "42"); !errors.Is(err, domain.ErrIntegrationNotFound) {
		t.Fatalf("owner mismatch = %v", err)
	}
}

func TestGetRepositorySupportsOwnerNameAndNormalizesProviderFailures(t *testing.T) {
	t.Parallel()
	status := http.StatusOK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/Unknowns24/akritas" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if status != http.StatusOK {
			http.Error(w, "provider detail must not leak", status)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 42, "name": "akritas", "default_branch": "main", "private": true, "html_url": "https://github.com/Unknowns24/akritas", "owner": map[string]any{"login": "Unknowns24"}})
	}))
	defer server.Close()
	account := githubAccount(t)
	credentials := credentialStoreFake{values: map[string][]byte{credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("secret")}}
	client := newTestClient(t, server.URL, credentials)
	if repository, err := client.GetRepository(context.Background(), account, "Unknowns24/akritas"); err != nil || repository.RepositoryIdentifier != "42" {
		t.Fatalf("owner/name resolution = %+v, %v", repository, err)
	}
	status = http.StatusNotFound
	if _, err := client.GetRepository(context.Background(), account, "Unknowns24/akritas"); !errors.Is(err, domain.ErrIntegrationNotFound) {
		t.Fatalf("not found = %v", err)
	}
	status = http.StatusUnauthorized
	if _, err := client.GetRepository(context.Background(), account, "Unknowns24/akritas"); !errors.Is(err, domain.ErrGitHubCredentialRejected) || strings.Contains(err.Error(), "provider detail") {
		t.Fatalf("credential rejection = %v", err)
	}
}

func newTestClient(t *testing.T, rawURL string, credentials portsout.CredentialStore) *Client {
	t.Helper()
	base, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(ClientConfig{APIBaseURL: base, HTTPClient: &http.Client{Timeout: time.Second}, Credentials: credentials, Now: time.Now})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func githubAccount(t *testing.T) domain.GitHubAccount {
	t.Helper()
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountPersonal, domain.GitHubAuthenticationPersonalAccessToken, "Unknowns24", domain.IntegrationStatusConnected, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	account.CredentialConfigured = true
	return *account
}

type credentialStoreFake struct{ values map[string][]byte }

func credentialKey(id uuid.UUID, kind portsout.SecretKind) string {
	return id.String() + ":" + string(kind)
}
func (f credentialStoreFake) Put(context.Context, string, uuid.UUID, portsout.SecretValue) error {
	return nil
}
func (f credentialStoreFake) Get(_ context.Context, _ string, id uuid.UUID, kind portsout.SecretKind) ([]byte, error) {
	value := f.values[credentialKey(id, kind)]
	copyValue := append([]byte(nil), value...)
	return copyValue, nil
}
func (f credentialStoreFake) DeleteOwner(context.Context, string, uuid.UUID) error { return nil }
func (f credentialStoreFake) MoveOwner(context.Context, string, uuid.UUID, string, uuid.UUID) error {
	return nil
}

var _ portsout.CredentialStore = credentialStoreFake{}

type bindingReaderFake struct{ binding portsout.GitHubAppBinding }

func (f bindingReaderFake) GetBinding(context.Context, uuid.UUID) (portsout.GitHubAppBinding, error) {
	return f.binding, nil
}
