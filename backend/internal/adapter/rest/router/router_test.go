package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	resthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler"
	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	dokployhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/dokploy"
	githubhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/github"
	projecthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/project"
	restmiddleware "github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type fakeAuthenticateSession struct {
	calls   int
	session domain.AdministratorSession
	err     error
}

func (f *fakeAuthenticateSession) Execute(context.Context, string) (domain.AdministratorSession, error) {
	f.calls++
	return f.session, f.err
}

type fakeGitHubAccounts struct {
	createCalls int
	getID       uuid.UUID
	panicOnGet  bool
}

func (f *fakeGitHubAccounts) CreatePAT(context.Context, portsin.CreateGitHubPATAccountCommand) (*domain.GitHubAccount, error) {
	f.createCalls++
	return &domain.GitHubAccount{}, nil
}

func (f *fakeGitHubAccounts) Get(_ context.Context, id uuid.UUID) (*domain.GitHubAccount, error) {
	if f.panicOnGet {
		panic("secret-panic-value")
	}
	f.getID = id
	return &domain.GitHubAccount{ID: id}, nil
}

func (f *fakeGitHubAccounts) List(context.Context, paging.Params) (paging.Slice[domain.GitHubAccount], error) {
	return paging.Slice[domain.GitHubAccount]{}, nil
}

func (f *fakeGitHubAccounts) Update(context.Context, uuid.UUID, portsin.UpdateGitHubAccountCommand) (*domain.GitHubAccount, error) {
	return &domain.GitHubAccount{}, nil
}

func (f *fakeGitHubAccounts) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeGitHubAccounts) TestConnection(context.Context, uuid.UUID) (portsin.ConnectionTestResult, error) {
	return portsin.ConnectionTestResult{}, nil
}

func (f *fakeGitHubAccounts) ListRepositories(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.GitHubRepository], error) {
	return paging.Slice[domain.GitHubRepository]{}, nil
}

type fakeGitHubApps struct {
	completeManifestCalls int
}

func (f *fakeGitHubApps) StartRegistration(context.Context, portsin.StartGitHubAppRegistrationCommand) (portsin.GitHubManifestRegistrationResult, error) {
	return portsin.GitHubManifestRegistrationResult{}, nil
}

func (f *fakeGitHubApps) CompleteManifest(context.Context, string, string) (portsin.GitHubManifestCallbackResult, error) {
	f.completeManifestCalls++
	return portsin.GitHubManifestCallbackResult{RedirectURL: "https://github.example/install"}, nil
}

func (f *fakeGitHubApps) CompleteInstallation(context.Context, int64, string) (portsin.GitHubInstallationCallbackResult, error) {
	return portsin.GitHubInstallationCallbackResult{RedirectURL: "https://app.example.com/settings"}, nil
}

type fakeDokployServers struct {
	getID uuid.UUID
}

type fakeProjects struct {
	createCalls int
}

func (f *fakeProjects) Create(context.Context, portsin.CreateProjectCommand) (*portsin.ProjectResult, error) {
	f.createCalls++
	return &portsin.ProjectResult{}, nil
}
func (*fakeProjects) Get(context.Context, uuid.UUID) (*portsin.ProjectResult, error) {
	return &portsin.ProjectResult{}, nil
}
func (*fakeProjects) List(context.Context, paging.Params) (paging.Slice[domain.Project], error) {
	return paging.Slice[domain.Project]{}, nil
}
func (*fakeProjects) Update(context.Context, portsin.UpdateProjectCommand) (*portsin.ProjectResult, error) {
	return &portsin.ProjectResult{}, nil
}
func (*fakeProjects) Delete(context.Context, uuid.UUID) error { return nil }
func (*fakeProjects) GetMonitoring(context.Context, uuid.UUID) (domain.MonitoringConfiguration, error) {
	return domain.MonitoringConfiguration{}, nil
}
func (*fakeProjects) PutMonitoring(context.Context, uuid.UUID, domain.MonitoringConfiguration) (domain.MonitoringConfiguration, error) {
	return domain.MonitoringConfiguration{}, nil
}

func (f *fakeDokployServers) Create(context.Context, portsin.CreateDokployServerCommand) (*domain.DokployServer, error) {
	return &domain.DokployServer{}, nil
}

func (f *fakeDokployServers) Get(_ context.Context, id uuid.UUID) (*domain.DokployServer, error) {
	f.getID = id
	return &domain.DokployServer{ID: id}, nil
}

func (f *fakeDokployServers) List(context.Context, paging.Params) (paging.Slice[domain.DokployServer], error) {
	return paging.Slice[domain.DokployServer]{}, nil
}

func (f *fakeDokployServers) Update(context.Context, uuid.UUID, portsin.UpdateDokployServerCommand) (*domain.DokployServer, error) {
	return &domain.DokployServer{}, nil
}

func (f *fakeDokployServers) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeDokployServers) TestConnection(context.Context, uuid.UUID) (portsin.ConnectionTestResult, error) {
	return portsin.ConnectionTestResult{}, nil
}

func (f *fakeDokployServers) ListApplications(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.DokployApplication], error) {
	return paging.Slice[domain.DokployApplication]{}, nil
}

type routerFixture struct {
	config       Config
	authenticate *fakeAuthenticateSession
	accounts     *fakeGitHubAccounts
	apps         *fakeGitHubApps
	dokploy      *fakeDokployServers
	projects     *fakeProjects
}

func newRouterFixture() *routerFixture {
	authenticate := &fakeAuthenticateSession{session: domain.AdministratorSession{ID: uuid.New()}}
	accounts := &fakeGitHubAccounts{}
	apps := &fakeGitHubApps{}
	dokploy := &fakeDokployServers{}
	paging := pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour}
	githubHandler, err := githubhandler.New(accounts, apps, paging)
	if err != nil {
		panic(err)
	}
	dokployHandler, err := dokployhandler.New(dokploy, paging)
	if err != nil {
		panic(err)
	}
	projects := &fakeProjects{}
	projectHandler, err := projecthandler.New(projects, paging)
	if err != nil {
		panic(err)
	}
	return &routerFixture{
		config: Config{
			Handlers: &resthandler.Handlers{
				AuthHandler:    &authhandler.Handler{},
				GitHubHandler:  githubHandler,
				DokployHandler: dokployHandler,
				ProjectHandler: projectHandler,
			},
			Admin:          restmiddleware.RequireSession(authenticate),
			Authenticate:   authenticate,
			AllowedOrigins: []string{"https://app.example.com"},
		},
		authenticate: authenticate,
		accounts:     accounts,
		apps:         apps,
		dokploy:      dokploy,
		projects:     projects,
	}
}

func validConfig() Config { return newRouterFixture().config }

func mustRouter(t *testing.T, config Config) http.Handler {
	t.Helper()
	handler, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return handler
}

func serve(handler http.Handler, method, target string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(method, target, nil))
	return recorder
}

func TestRouterFailsClosedWithoutAdministratorMiddleware(t *testing.T) {
	handler, err := New(Config{})
	if !errors.Is(err, ErrAdminMiddlewareUnavailable) || handler != nil {
		t.Fatalf("router must not be mountable without administrator middleware: handler=%v err=%v", handler, err)
	}
}

func TestRouterFailsClosedWhenAdministratorMiddlewareProducesNil(t *testing.T) {
	config := validConfig()
	config.Admin = func(http.Handler) http.Handler { return nil }

	handler, err := New(config)
	if !errors.Is(err, ErrAdminMiddlewareUnavailable) || handler != nil {
		t.Fatalf("New() = (%v, %v), want fail-closed middleware error", handler, err)
	}
}

func TestRouterFailsClosedWithoutCompleteHandlers(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config)
	}{
		{name: "missing aggregate", mutate: func(config *Config) { config.Handlers = nil }},
		{name: "missing auth", mutate: func(config *Config) { config.Handlers.AuthHandler = nil }},
		{name: "missing GitHub", mutate: func(config *Config) { config.Handlers.GitHubHandler = nil }},
		{name: "missing Dokploy", mutate: func(config *Config) { config.Handlers.DokployHandler = nil }},
		{name: "missing Project", mutate: func(config *Config) { config.Handlers.ProjectHandler = nil }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := validConfig()
			test.mutate(&config)
			handler, err := New(config)
			if handler != nil || !errors.Is(err, ErrInvalidRouterConfiguration) {
				t.Fatalf("New() = (%v, %v), want invalid configuration", handler, err)
			}
		})
	}
}

func TestRouterExposesExactChiRouteInventory(t *testing.T) {
	handler := mustRouter(t, validConfig())
	routes, ok := handler.(chi.Routes)
	if !ok {
		t.Fatalf("router type %T does not expose chi routes", handler)
	}

	var got []string
	if err := chi.Walk(routes, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		got = append(got, method+" "+route)
		return nil
	}); err != nil {
		t.Fatalf("walk routes: %v", err)
	}
	sort.Strings(got)

	want := []string{
		"DELETE /api/v1/auth/session",
		"DELETE /api/v1/integrations/dokploy/servers/{server_id}",
		"DELETE /api/v1/integrations/github/accounts/{account_id}",
		"DELETE /api/v1/projects/{project_id}",
		"GET /api/v1/auth/session",
		"GET /api/v1/auth/setup-status",
		"GET /api/v1/integrations/dokploy/servers",
		"GET /api/v1/integrations/dokploy/servers/{server_id}",
		"GET /api/v1/integrations/dokploy/servers/{server_id}/applications",
		"GET /api/v1/integrations/github/accounts",
		"GET /api/v1/integrations/github/accounts/{account_id}",
		"GET /api/v1/integrations/github/accounts/{account_id}/repositories",
		"GET /api/v1/integrations/github/app-installations/callback",
		"GET /api/v1/integrations/github/app-manifest/callback",
		"GET /api/v1/projects",
		"GET /api/v1/projects/{project_id}",
		"GET /api/v1/projects/{project_id}/monitoring-configuration",
		"PATCH /api/v1/integrations/dokploy/servers/{server_id}",
		"PATCH /api/v1/integrations/github/accounts/{account_id}",
		"PATCH /api/v1/projects/{project_id}",
		"POST /api/v1/auth/login",
		"POST /api/v1/auth/setup",
		"POST /api/v1/auth/setup/verify",
		"POST /api/v1/integrations/dokploy/servers",
		"POST /api/v1/integrations/dokploy/servers/{server_id}/connection-test",
		"POST /api/v1/integrations/github/accounts",
		"POST /api/v1/integrations/github/accounts/{account_id}/connection-test",
		"POST /api/v1/integrations/github/app-manifest/registrations",
		"POST /api/v1/projects",
		"PUT /api/v1/projects/{project_id}/monitoring-configuration",
	}
	sort.Strings(want)
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("route inventory mismatch\ngot:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestRouterKeepsGitHubCallbacksPublicAndSupportsHead(t *testing.T) {
	fixture := newRouterFixture()
	fixture.authenticate.err = domain.ErrInactiveAdministratorSession
	handler := mustRouter(t, fixture.config)

	for _, method := range []string{http.MethodGet, http.MethodHead} {
		recorder := serve(handler, method, "/api/v1/integrations/github/app-manifest/callback?code=code&state=state")
		if recorder.Code != http.StatusSeeOther {
			t.Fatalf("%s callback status = %d, want 303", method, recorder.Code)
		}
	}
	if fixture.authenticate.calls != 0 {
		t.Fatalf("public callback authenticated %d times", fixture.authenticate.calls)
	}
	if fixture.apps.completeManifestCalls != 2 {
		t.Fatalf("callback calls = %d, want 2", fixture.apps.completeManifestCalls)
	}
}

func TestRouterProtectsSessionAndIntegrationMutations(t *testing.T) {
	t.Run("session requires authentication", func(t *testing.T) {
		fixture := newRouterFixture()
		fixture.authenticate.err = domain.ErrInactiveAdministratorSession
		recorder := serve(mustRouter(t, fixture.config), http.MethodGet, "/api/v1/auth/session")
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", recorder.Code)
		}
	})

	t.Run("logout checks origin after authentication", func(t *testing.T) {
		fixture := newRouterFixture()
		recorder := serve(mustRouter(t, fixture.config), http.MethodDelete, "/api/v1/auth/session")
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403", recorder.Code)
		}
		if fixture.authenticate.calls != 1 {
			t.Fatalf("authentication calls = %d, want 1", fixture.authenticate.calls)
		}
	})

	t.Run("private mutation checks origin after admin", func(t *testing.T) {
		fixture := newRouterFixture()
		recorder := serve(mustRouter(t, fixture.config), http.MethodPost, "/api/v1/integrations/github/accounts")
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403", recorder.Code)
		}
		if fixture.authenticate.calls != 1 {
			t.Fatalf("authentication calls = %d, want 1", fixture.authenticate.calls)
		}
		if fixture.accounts.createCalls != 0 {
			t.Fatal("handler must not run when Origin is rejected")
		}
	})

	t.Run("Project requires authentication", func(t *testing.T) {
		fixture := newRouterFixture()
		fixture.authenticate.err = domain.ErrInactiveAdministratorSession
		recorder := serve(mustRouter(t, fixture.config), http.MethodGet, "/api/v1/projects")
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", recorder.Code)
		}
	})

	t.Run("Project mutation requires allowed Origin", func(t *testing.T) {
		fixture := newRouterFixture()
		recorder := serve(mustRouter(t, fixture.config), http.MethodPost, "/api/v1/projects")
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403", recorder.Code)
		}
		if fixture.projects.createCalls != 0 {
			t.Fatal("Project handler must not run when Origin is rejected")
		}
	})
}

func TestRouterPopulatesStandardLibraryPathValues(t *testing.T) {
	fixture := newRouterFixture()
	handler := mustRouter(t, fixture.config)
	accountID := uuid.New()
	serverID := uuid.New()

	if recorder := serve(handler, http.MethodGet, "/api/v1/integrations/github/accounts/"+accountID.String()); recorder.Code != http.StatusOK {
		t.Fatalf("GitHub status = %d, want 200", recorder.Code)
	}
	if recorder := serve(handler, http.MethodGet, "/api/v1/integrations/dokploy/servers/"+serverID.String()); recorder.Code != http.StatusOK {
		t.Fatalf("Dokploy status = %d, want 200", recorder.Code)
	}
	if fixture.accounts.getID != accountID || fixture.dokploy.getID != serverID {
		t.Fatalf("path values = (%s, %s), want (%s, %s)", fixture.accounts.getID, fixture.dokploy.getID, accountID, serverID)
	}
}

func TestRouterUsesExactPathsAndStandardRoutingErrors(t *testing.T) {
	handler := mustRouter(t, validConfig())

	if recorder := serve(handler, http.MethodGet, "/api/v1/unknown"); recorder.Code != http.StatusNotFound {
		t.Fatalf("unknown status = %d, want 404", recorder.Code)
	}
	if recorder := serve(handler, http.MethodGet, "/api/v1/auth/setup-status/"); recorder.Code != http.StatusNotFound {
		t.Fatalf("trailing slash status = %d, want 404", recorder.Code)
	}
	recorder := serve(handler, http.MethodPut, "/api/v1/auth/setup-status")
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("wrong method status = %d, want 405", recorder.Code)
	}
	if !strings.Contains(recorder.Header().Get("Allow"), http.MethodGet) {
		t.Fatalf("Allow = %q, want GET", recorder.Header().Get("Allow"))
	}
}

func TestRouterRecoversPanicsWithStableJSONAndRequestID(t *testing.T) {
	for _, requestID := range []string{"req-client-12345", "", "bad"} {
		t.Run("request_id="+requestID, func(t *testing.T) {
			fixture := newRouterFixture()
			fixture.accounts.panicOnGet = true
			handler := mustRouter(t, fixture.config)
			request := httptest.NewRequest(http.MethodGet, "/api/v1/integrations/github/accounts/"+uuid.NewString(), nil)
			if requestID != "" {
				request.Header.Set("X-Request-ID", requestID)
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want 500", recorder.Code)
			}
			body := recorder.Body.String()
			if !strings.Contains(body, `"code":"0x102002I"`) {
				t.Fatalf("missing stable error code: %s", body)
			}
			if strings.Contains(body, "secret-panic-value") {
				t.Fatal("panic value leaked in response")
			}
			if requestID == "req-client-12345" && !strings.Contains(body, `"request_id":"req-client-12345"`) {
				t.Fatalf("valid request ID not preserved: %s", body)
			}
			if requestID == "bad" && strings.Contains(body, `"request_id":"bad"`) {
				t.Fatalf("invalid request ID was trusted: %s", body)
			}
		})
	}
}
