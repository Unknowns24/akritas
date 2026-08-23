package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func TestCallbacksAreNoStoreAndOnlyUseBackendSelectedRedirects(t *testing.T) {
	pagingConfig, _ := pagination.NewConfig([]byte("01234567890123456789012345678901"), time.Minute)
	apps := appUseCaseStub{manifestRedirect: "https://github.com/apps/akritas/installations/new?state=server-state", installationRedirect: "https://akritas.example.com/settings/integrations/github"}
	handler, err := New(accountUseCaseStub{}, apps, pagingConfig)
	if err != nil {
		t.Fatal(err)
	}
	manifestRequest := httptest.NewRequest(http.MethodGet, "/api/v1/integrations/github/app-manifest/callback?code=manifest-code-at-least-twenty&state=browser-supplied-state-that-is-long-enough&redirect=https://attacker.example", nil)
	manifestResponse := httptest.NewRecorder()
	handler.CompleteManifest(manifestResponse, manifestRequest)
	if manifestResponse.Code != http.StatusSeeOther || manifestResponse.Header().Get("Cache-Control") != "no-store" || manifestResponse.Header().Get("Location") != apps.manifestRedirect {
		t.Fatalf("unsafe manifest callback response: code=%d headers=%v", manifestResponse.Code, manifestResponse.Header())
	}
	installationRequest := httptest.NewRequest(http.MethodGet, "/api/v1/integrations/github/app-installations/callback?installation_id=99&state=browser-supplied-state-that-is-long-enough&redirect=https://attacker.example", nil)
	installationResponse := httptest.NewRecorder()
	handler.CompleteInstallation(installationResponse, installationRequest)
	if installationResponse.Code != http.StatusSeeOther || installationResponse.Header().Get("Cache-Control") != "no-store" || installationResponse.Header().Get("Location") != apps.installationRedirect {
		t.Fatalf("unsafe installation callback response: code=%d headers=%v", installationResponse.Code, installationResponse.Header())
	}
}

type appUseCaseStub struct {
	manifestRedirect     string
	installationRedirect string
}

func (appUseCaseStub) StartRegistration(context.Context, portsin.StartGitHubAppRegistrationCommand) (portsin.GitHubManifestRegistrationResult, error) {
	return portsin.GitHubManifestRegistrationResult{}, nil
}
func (s appUseCaseStub) CompleteManifest(context.Context, string, string) (portsin.GitHubManifestCallbackResult, error) {
	return portsin.GitHubManifestCallbackResult{RedirectURL: s.manifestRedirect}, nil
}
func (s appUseCaseStub) CompleteInstallation(context.Context, int64, string) (portsin.GitHubInstallationCallbackResult, error) {
	return portsin.GitHubInstallationCallbackResult{RedirectURL: s.installationRedirect}, nil
}

type accountUseCaseStub struct{}

func (accountUseCaseStub) CreatePAT(context.Context, portsin.CreateGitHubPATAccountCommand) (*domain.GitHubAccount, error) {
	return nil, nil
}
func (accountUseCaseStub) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return nil, nil
}
func (accountUseCaseStub) List(context.Context, paging.Params) (paging.Slice[domain.GitHubAccount], error) {
	return paging.Slice[domain.GitHubAccount]{}, nil
}
func (accountUseCaseStub) Update(context.Context, uuid.UUID, portsin.UpdateGitHubAccountCommand) (*domain.GitHubAccount, error) {
	return nil, nil
}
func (accountUseCaseStub) Delete(context.Context, uuid.UUID) error { return nil }
func (accountUseCaseStub) TestConnection(context.Context, uuid.UUID) (portsin.ConnectionTestResult, error) {
	return portsin.ConnectionTestResult{}, nil
}
func (accountUseCaseStub) ListRepositories(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.GitHubRepository], error) {
	return paging.Slice[domain.GitHubRepository]{}, nil
}
