package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/external/integration"
	projecthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	projectuc "github.com/Unknowns24/akritas/backend/internal/usecase/project"
	"github.com/google/uuid"
)

func TestProjectHTTPContract(t *testing.T) {
	handler, accountID, serverID := newAPI(t)

	health := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	healthRec := httptest.NewRecorder()
	handler.ServeHTTP(healthRec, health)
	if healthRec.Code != http.StatusOK {
		t.Fatalf("health: %d %s", healthRec.Code, healthRec.Body.Bytes())
	}

	unauthorized := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	unauthRec := httptest.NewRecorder()
	handler.ServeHTTP(unauthRec, unauthorized)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d %s", unauthRec.Code, unauthRec.Body.Bytes())
	}

	body := fmt.Sprintf(`{
		"name":"sentinel-api",
		"description":"demo",
		"github_account_id":%q,
		"repository_identifier":"Unknowns24/akritas",
		"default_branch":"main",
		"dokploy_server_id":%q,
		"application_identifier":"app-1",
		"monitoring_configuration":{
			"enabled":false,
			"error_patterns":[],
			"ignored_patterns":[],
			"grouping_window":"PT30M",
			"context_before":20,
			"context_after":20
		}
	}`, accountID, serverID)
	createRec := doJSON(t, handler, http.MethodPost, "/api/v1/projects", body, true)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", createRec.Code, createRec.Body.Bytes())
	}
	assertNoSecretFields(t, createRec.Body.Bytes())
	var created struct {
		Data struct {
			ID                      string `json:"id"`
			MonitoringStatus        string `json:"monitoring_status"`
			BuiltInDetectionRules   []any  `json:"built_in_detection_rules"`
			MonitoringConfiguration struct {
				GroupingWindow string `json:"grouping_window"`
			} `json:"monitoring_configuration"`
			GitHubRepository struct {
				Owner string `json:"owner"`
				Name  string `json:"name"`
			} `json:"github_repository"`
		} `json:"data"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.Data.MonitoringStatus != "disabled" || len(created.Data.BuiltInDetectionRules) != 7 {
		t.Fatalf("create payload: %+v", created.Data)
	}
	if created.Data.GitHubRepository.Owner != "Unknowns24" || created.Data.MonitoringConfiguration.GroupingWindow != "PT30M" {
		t.Fatalf("snapshot/duration: %+v", created.Data)
	}

	getRec := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID, "", true)
	if getRec.Code != http.StatusOK {
		t.Fatalf("get: %d %s", getRec.Code, getRec.Body.Bytes())
	}

	listRec := doJSON(t, handler, http.MethodGet, "/api/v1/projects?limit=10&name_like=sentinel", "", true)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list: %d %s", listRec.Code, listRec.Body.Bytes())
	}
	if !bytes.Contains(listRec.Body.Bytes(), []byte(`"paging"`)) {
		t.Fatalf("list missing paging: %s", listRec.Body.Bytes())
	}

	missingAccount := strings.ReplaceAll(body, accountID.String(), uuid.NewString())
	missingRec := doJSON(t, handler, http.MethodPost, "/api/v1/projects", missingAccount, true)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing account: %d %s", missingRec.Code, missingRec.Body.Bytes())
	}

	dupName := strings.ReplaceAll(body, `"app-1"`, `"app-2"`)
	dupRec := doJSON(t, handler, http.MethodPost, "/api/v1/projects", dupName, true)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("duplicate name: %d %s", dupRec.Code, dupRec.Body.Bytes())
	}

	dupApp := strings.ReplaceAll(body, `"name":"sentinel-api"`, `"name":"other-api"`)
	dupAppRec := doJSON(t, handler, http.MethodPost, "/api/v1/projects", dupApp, true)
	if dupAppRec.Code != http.StatusConflict {
		t.Fatalf("duplicate application: %d %s", dupAppRec.Code, dupAppRec.Body.Bytes())
	}

	invalid := strings.ReplaceAll(body, `"name":"sentinel-api"`, `"name":"other"`)
	invalid = strings.ReplaceAll(invalid, `"grouping_window":"PT30M"`, `"grouping_window":"nope"`)
	invalidRec := doJSON(t, handler, http.MethodPost, "/api/v1/projects", invalid, true)
	if invalidRec.Code != http.StatusBadRequest {
		t.Fatalf("invalid duration: %d %s", invalidRec.Code, invalidRec.Body.Bytes())
	}

	enable := `{
		"enabled":true,
		"error_patterns":["panic"],
		"ignored_patterns":["healthcheck"],
		"grouping_window":"PT15M",
		"context_before":5,
		"context_after":5
	}`
	putRec := doJSON(t, handler, http.MethodPut, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", enable, true)
	if putRec.Code != http.StatusOK {
		t.Fatalf("put monitoring: %d %s", putRec.Code, putRec.Body.Bytes())
	}
	var putBody struct {
		Data struct {
			Enabled         bool     `json:"enabled"`
			ErrorPatterns   []string `json:"error_patterns"`
			IgnoredPatterns []string `json:"ignored_patterns"`
			GroupingWindow  string   `json:"grouping_window"`
			ContextBefore   int      `json:"context_before"`
			ContextAfter    int      `json:"context_after"`
		} `json:"data"`
	}
	if err := json.Unmarshal(putRec.Body.Bytes(), &putBody); err != nil {
		t.Fatal(err)
	}
	if !putBody.Data.Enabled || putBody.Data.GroupingWindow != "PT15M" ||
		putBody.Data.ContextBefore != 5 || putBody.Data.ContextAfter != 5 ||
		len(putBody.Data.ErrorPatterns) != 1 || putBody.Data.ErrorPatterns[0] != "panic" ||
		len(putBody.Data.IgnoredPatterns) != 1 || putBody.Data.IgnoredPatterns[0] != "healthcheck" {
		t.Fatalf("put envelope: %+v", putBody.Data)
	}

	getConfig := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", "", true)
	if getConfig.Code != http.StatusOK {
		t.Fatalf("get monitoring: %d %s", getConfig.Code, getConfig.Body.Bytes())
	}
	if err := json.Unmarshal(getConfig.Body.Bytes(), &putBody); err != nil {
		t.Fatal(err)
	}
	if !putBody.Data.Enabled || putBody.Data.GroupingWindow != "PT15M" ||
		putBody.Data.ContextBefore != 5 || putBody.Data.ContextAfter != 5 {
		t.Fatalf("get monitoring envelope: %+v", putBody.Data)
	}

	after := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID, "", true)
	if !bytes.Contains(after.Body.Bytes(), []byte(`"monitoring_status":"starting"`)) {
		t.Fatalf("expected starting: %s", after.Body.Bytes())
	}

	badRegex := `{
		"enabled":false,
		"error_patterns":["["],
		"ignored_patterns":[],
		"grouping_window":"PT30M",
		"context_before":20,
		"context_after":20
	}`
	badRegexRec := doJSON(t, handler, http.MethodPut, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", badRegex, true)
	if badRegexRec.Code != http.StatusBadRequest || !bytes.Contains(badRegexRec.Body.Bytes(), []byte(`"0x403004V"`)) {
		t.Fatalf("invalid regex: %d %s", badRegexRec.Code, badRegexRec.Body.Bytes())
	}

	badDuration := strings.ReplaceAll(enable, `"PT15M"`, `"nope"`)
	badDurationRec := doJSON(t, handler, http.MethodPut, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", badDuration, true)
	if badDurationRec.Code != http.StatusBadRequest || !bytes.Contains(badDurationRec.Body.Bytes(), []byte(`"0x403004V"`)) {
		t.Fatalf("invalid duration: %d %s", badDurationRec.Code, badDurationRec.Body.Bytes())
	}

	missingConfig := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+uuid.NewString()+"/monitoring-configuration", "", true)
	if missingConfig.Code != http.StatusNotFound {
		t.Fatalf("missing project config: %d %s", missingConfig.Code, missingConfig.Body.Bytes())
	}
	unauthConfig := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", "", false)
	if unauthConfig.Code != http.StatusUnauthorized {
		t.Fatalf("unauth config: %d %s", unauthConfig.Code, unauthConfig.Body.Bytes())
	}

	notFound := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+uuid.NewString(), "", true)
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("missing project: %d", notFound.Code)
	}
}

func TestProjectHTTPEnableDisableSnapshot(t *testing.T) {
	handler, accountID, serverID := newAPI(t)

	body := fmt.Sprintf(`{
		"name":"sentinel-api",
		"github_account_id":%q,
		"repository_identifier":"Unknowns24/akritas",
		"default_branch":"main",
		"dokploy_server_id":%q,
		"application_identifier":"app-1",
		"monitoring_configuration":{
			"enabled":false,
			"error_patterns":[],
			"ignored_patterns":[],
			"grouping_window":"PT30M",
			"context_before":20,
			"context_after":20
		}
	}`, accountID, serverID)
	createRec := doJSON(t, handler, http.MethodPost, "/api/v1/projects", body, true)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", createRec.Code, createRec.Body.Bytes())
	}
	var created struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	unauth := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID, "", false)
	if unauth.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without cookie, got %d %s", unauth.Code, unauth.Body.Bytes())
	}

	enable := `{
		"enabled":true,
		"error_patterns":["panic"],
		"ignored_patterns":["healthcheck"],
		"grouping_window":"PT15M",
		"context_before":8,
		"context_after":12
	}`
	putRec := doJSON(t, handler, http.MethodPut, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", enable, true)
	if putRec.Code != http.StatusOK {
		t.Fatalf("enable: %d %s", putRec.Code, putRec.Body.Bytes())
	}

	enabledRec := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID, "", true)
	if enabledRec.Code != http.StatusOK {
		t.Fatalf("get enabled: %d %s", enabledRec.Code, enabledRec.Body.Bytes())
	}
	assertNoSecretFields(t, enabledRec.Body.Bytes())
	got := decodeProjectEnvelope(t, enabledRec.Body.Bytes())
	if got.MonitoringStatus != "starting" {
		t.Fatalf("expected starting, got %s", got.MonitoringStatus)
	}
	repo := got.GitHubRepository
	if repo.GitHubAccountID != accountID.String() || repo.RepositoryIdentifier != "Unknowns24/akritas" {
		t.Fatalf("github identity: %+v", repo)
	}
	if repo.Owner != "Unknowns24" || repo.Name != "akritas" || repo.FullName != "Unknowns24/akritas" {
		t.Fatalf("github names: %+v", repo)
	}
	if repo.DefaultBranch != "main" || repo.Private || repo.HTMLURL != "https://github.com/Unknowns24/akritas" {
		t.Fatalf("github snapshot: %+v", repo)
	}
	app := got.DokployApplication
	if app.DokployServerID != serverID.String() || app.ApplicationIdentifier != "app-1" {
		t.Fatalf("dokploy identity: %+v", app)
	}
	if app.InstanceIdentifier != "app-1" || app.DisplayName != "app-1" || app.Status != "unknown" {
		t.Fatalf("dokploy snapshot: %+v", app)
	}
	cfg := got.MonitoringConfiguration
	if !cfg.Enabled || cfg.GroupingWindow != "PT15M" || cfg.ContextBefore != 8 || cfg.ContextAfter != 12 {
		t.Fatalf("enabled config: %+v", cfg)
	}
	if len(cfg.ErrorPatterns) != 1 || cfg.ErrorPatterns[0] != "panic" ||
		len(cfg.IgnoredPatterns) != 1 || cfg.IgnoredPatterns[0] != "healthcheck" {
		t.Fatalf("enabled patterns: %+v", cfg)
	}

	disable := `{
		"enabled":false,
		"error_patterns":["panic"],
		"ignored_patterns":["healthcheck"],
		"grouping_window":"PT15M",
		"context_before":8,
		"context_after":12
	}`
	disableRec := doJSON(t, handler, http.MethodPut, "/api/v1/projects/"+created.Data.ID+"/monitoring-configuration", disable, true)
	if disableRec.Code != http.StatusOK {
		t.Fatalf("disable: %d %s", disableRec.Code, disableRec.Body.Bytes())
	}

	disabledRec := doJSON(t, handler, http.MethodGet, "/api/v1/projects/"+created.Data.ID, "", true)
	if disabledRec.Code != http.StatusOK {
		t.Fatalf("get disabled: %d %s", disabledRec.Code, disabledRec.Body.Bytes())
	}
	assertNoSecretFields(t, disabledRec.Body.Bytes())
	off := decodeProjectEnvelope(t, disabledRec.Body.Bytes())
	if off.MonitoringStatus != "disabled" {
		t.Fatalf("expected disabled, got %s", off.MonitoringStatus)
	}
	if off.GitHubRepository != repo || off.DokployApplication != app {
		t.Fatalf("disable must keep snapshots: github=%+v dokploy=%+v", off.GitHubRepository, off.DokployApplication)
	}
	if off.MonitoringConfiguration.Enabled || off.MonitoringConfiguration.GroupingWindow != "PT15M" ||
		off.MonitoringConfiguration.ContextBefore != 8 || off.MonitoringConfiguration.ContextAfter != 12 {
		t.Fatalf("disable must keep config: %+v", off.MonitoringConfiguration)
	}
	if len(off.MonitoringConfiguration.ErrorPatterns) != 1 || off.MonitoringConfiguration.ErrorPatterns[0] != "panic" ||
		len(off.MonitoringConfiguration.IgnoredPatterns) != 1 || off.MonitoringConfiguration.IgnoredPatterns[0] != "healthcheck" {
		t.Fatalf("disable must keep patterns: %+v", off.MonitoringConfiguration)
	}
}

type projectHTTPEnvelope struct {
	MonitoringStatus string `json:"monitoring_status"`
	GitHubRepository struct {
		GitHubAccountID      string `json:"github_account_id"`
		RepositoryIdentifier string `json:"repository_identifier"`
		Owner                string `json:"owner"`
		Name                 string `json:"name"`
		FullName             string `json:"full_name"`
		DefaultBranch        string `json:"default_branch"`
		Private              bool   `json:"private"`
		HTMLURL              string `json:"html_url"`
	} `json:"github_repository"`
	DokployApplication struct {
		DokployServerID       string `json:"dokploy_server_id"`
		ApplicationIdentifier string `json:"application_identifier"`
		InstanceIdentifier    string `json:"instance_identifier"`
		DisplayName           string `json:"display_name"`
		Environment           string `json:"environment"`
		Status                string `json:"status"`
	} `json:"dokploy_application"`
	MonitoringConfiguration struct {
		Enabled         bool     `json:"enabled"`
		ErrorPatterns   []string `json:"error_patterns"`
		IgnoredPatterns []string `json:"ignored_patterns"`
		GroupingWindow  string   `json:"grouping_window"`
		ContextBefore   int      `json:"context_before"`
		ContextAfter    int      `json:"context_after"`
	} `json:"monitoring_configuration"`
}

func decodeProjectEnvelope(t *testing.T, payload []byte) projectHTTPEnvelope {
	t.Helper()
	var envelope struct {
		Data projectHTTPEnvelope `json:"data"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope.Data
}

func newAPI(t *testing.T) (http.Handler, uuid.UUID, uuid.UUID) {
	t.Helper()
	db := dbtest.OpenMigrated(t)
	now := time.Date(2026, 8, 22, 18, 0, 0, 0, time.UTC)
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationGitHubApp, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "demo", "https://dokploy.example.com", "server-1", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := githubaccount.NewRepository(db).Create(context.Background(), account); err != nil {
		t.Fatal(err)
	}
	if err := dokployserver.NewRepository(db).Create(context.Background(), server); err != nil {
		t.Fatal(err)
	}
	useCase := projectuc.NewUseCase(
		projectrepo.NewRepository(db),
		githubaccount.NewRepository(db),
		dokployserver.NewRepository(db),
		integration.NewSnapshotResolver(),
	)
	projects := projecthandler.NewHandler(useCase, useCase, useCase, useCase, useCase, useCase, "test-secret")
	return router.New(projects, middleware.AllowNonEmptySessions{}), account.ID, server.ID
}

func doJSON(t *testing.T, handler http.Handler, method, path, body string, auth bool) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body == "" {
		reader = bytes.NewReader(nil)
	} else {
		reader = bytes.NewReader([]byte(body))
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth {
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookie, Value: "test-session"})
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func assertNoSecretFields(t *testing.T, payload []byte) {
	t.Helper()
	encoded := strings.ToLower(string(payload))
	for _, leaked := range []string{`"token"`, `"password"`, `"api_key"`, `"api_credential"`, `"pat"`, `"master_key"`, `"secret"`} {
		if strings.Contains(encoded, leaked) {
			t.Fatalf("payload leaked %s: %s", leaked, payload)
		}
	}
}
