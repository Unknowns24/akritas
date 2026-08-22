package project

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func TestCreateGetListUpdateAndMonitoring(t *testing.T) {
	t.Parallel()

	uc, account, server := newTestUseCase(t)
	ctx := context.Background()
	created, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "sentinel-api", Description: "demo",
		GitHubAccountID: account.ID, RepositoryIdentifier: "Unknowns24/akritas", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-1",
		MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	assertNoSecrets(t, created.Project)
	if created.Project.GitHubRepository.Owner != "Unknowns24" || created.Project.GitHubRepository.Name != "akritas" {
		t.Fatalf("owner/name snapshot: %+v", created.Project.GitHubRepository)
	}
	if created.Project.MonitoringStatus != domain.MonitoringStatusDisabled {
		t.Fatalf("monitoring should start disabled, got %s", created.Project.MonitoringStatus)
	}
	if len(created.BuiltInDetectionRules) != 7 {
		t.Fatalf("expected built-in rules, got %d", len(created.BuiltInDetectionRules))
	}

	got, err := uc.Get(ctx, created.Project.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Project.Name != "sentinel-api" {
		t.Fatalf("get returned %+v", got.Project)
	}

	listed, total, err := uc.List(ctx, paging.ListQuery{Limit: 10, NameLike: "sentinel"})
	if err != nil || total != 1 || len(listed) != 1 {
		t.Fatalf("list: listed=%d total=%d err=%v", len(listed), total, err)
	}

	newName := "sentinel-core"
	updated, err := uc.Update(ctx, inproject.UpdateCommand{ID: created.Project.ID, Name: &newName})
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	if updated.Project.Name != "sentinel-core" {
		t.Fatalf("rename failed: %+v", updated.Project)
	}

	enabled, err := domain.NewMonitoringConfiguration(true, []string{`panic`}, []string{`healthcheck`}, time.Minute, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	config, err := uc.PutMonitoring(ctx, inproject.MonitoringCommand{ProjectID: created.Project.ID, MonitoringConfiguration: enabled})
	if err != nil {
		t.Fatalf("enable: %v", err)
	}
	if !config.Enabled {
		t.Fatal("monitoring was not enabled")
	}
	afterEnable, err := uc.Get(ctx, created.Project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if afterEnable.Project.MonitoringStatus != domain.MonitoringStatusStarting {
		t.Fatalf("expected starting, got %s", afterEnable.Project.MonitoringStatus)
	}

	afterEnable.Project.MonitoringStatus = domain.MonitoringStatusDegraded
	if err := uc.projects.Update(ctx, afterEnable.Project); err != nil {
		t.Fatal(err)
	}
	if _, err := uc.PutMonitoring(ctx, inproject.MonitoringCommand{ProjectID: created.Project.ID, MonitoringConfiguration: enabled}); err != nil {
		t.Fatal(err)
	}
	degraded, err := uc.Get(ctx, created.Project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if degraded.Project.MonitoringStatus != domain.MonitoringStatusDegraded {
		t.Fatalf("degraded status was reset to %s", degraded.Project.MonitoringStatus)
	}

	disabled := domain.DefaultMonitoringConfiguration()
	if _, err := uc.PutMonitoring(ctx, inproject.MonitoringCommand{ProjectID: created.Project.ID, MonitoringConfiguration: disabled}); err != nil {
		t.Fatal(err)
	}
	afterDisable, err := uc.GetMonitoring(ctx, created.Project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if afterDisable.Enabled {
		t.Fatal("monitoring should be disabled")
	}
}

func TestCreateRejectsMissingIntegrationsDuplicatesAndSecrets(t *testing.T) {
	t.Parallel()

	uc, account, server := newTestUseCase(t)
	ctx := context.Background()
	if _, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "one", GitHubAccountID: uuid.New(), RepositoryIdentifier: "Unknowns24/akritas", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-1", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	}); !errors.Is(err, apperr.ErrGitHubAccountNotFound) {
		t.Fatalf("expected missing account, got %v", err)
	}
	if _, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "one", GitHubAccountID: account.ID, RepositoryIdentifier: "Unknowns24/akritas", DefaultBranch: "main",
		DokployServerID: uuid.New(), ApplicationIdentifier: "app-1", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	}); !errors.Is(err, apperr.ErrDokployServerNotFound) {
		t.Fatalf("expected missing server, got %v", err)
	}

	first, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "sentinel-api", GitHubAccountID: account.ID, RepositoryIdentifier: "Unknowns24/akritas", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-1", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "SENTINEL-API", GitHubAccountID: account.ID, RepositoryIdentifier: "Unknowns24/other", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-2", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	}); !errors.Is(err, apperr.ErrProjectNameConflict) {
		t.Fatalf("expected name conflict, got %v", err)
	}
	if _, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "other", GitHubAccountID: account.ID, RepositoryIdentifier: "Unknowns24/other", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-1", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	}); !errors.Is(err, apperr.ErrProjectApplicationConflict) {
		t.Fatalf("expected application conflict, got %v", err)
	}
	assertNoSecrets(t, first.Project)
}

func TestCreateRejectsUnresolvableRepositoryIdentifier(t *testing.T) {
	t.Parallel()

	uc, account, server := newTestUseCase(t)
	ctx := context.Background()
	if _, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "broken", GitHubAccountID: account.ID, RepositoryIdentifier: "a/b/c", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-1", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	}); !errors.Is(err, apperr.ErrRepositoryNotResolvable) {
		t.Fatalf("expected unresolvable repository, got %v", err)
	}
}

func TestUpdatePersistsNewGitHubSnapshot(t *testing.T) {
	t.Parallel()

	uc, account, server := newTestUseCase(t)
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
		other, err := domain.NewGitHubAccount(uuid.New(), "Other", domain.GitHubAccountPersonal, domain.GitHubAuthenticationGitHubApp, "other-org", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	uc.accounts.(*memoryAccounts).byID[other.ID] = other

	ctx := context.Background()
	created, err := uc.Create(ctx, inproject.CreateCommand{
		Name: "sentinel-api", GitHubAccountID: account.ID, RepositoryIdentifier: "Unknowns24/akritas", DefaultBranch: "main",
		DokployServerID: server.ID, ApplicationIdentifier: "app-1", MonitoringConfiguration: domain.DefaultMonitoringConfiguration(),
	})
	if err != nil {
		t.Fatal(err)
	}

	identifier := "other-org/sentinel"
	branch := "develop"
	updated, err := uc.Update(ctx, inproject.UpdateCommand{
		ID: created.Project.ID, GitHubAccountID: &other.ID, RepositoryIdentifier: &identifier, DefaultBranch: &branch,
	})
	if err != nil {
		t.Fatalf("update integrations: %v", err)
	}
	repo := updated.Project.GitHubRepository
	if repo.GitHubAccountID != other.ID || repo.RepositoryIdentifier != identifier {
		t.Fatalf("identity not updated: %+v", repo)
	}
	if repo.Owner != "other-org" || repo.Name != "sentinel" || repo.FullName != "other-org/sentinel" {
		t.Fatalf("snapshot not rebuilt: %+v", repo)
	}
	if repo.DefaultBranch != "develop" || repo.HTMLURL != "https://github.com/other-org/sentinel" || repo.Private {
		t.Fatalf("projection mismatch: %+v", repo)
	}

	got, err := uc.Get(ctx, created.Project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Project.GitHubRepository != repo {
		t.Fatalf("persisted snapshot mismatch: %+v", got.Project.GitHubRepository)
	}
	assertNoSecrets(t, got.Project)

	count, err := uc.projects.CountByGitHubAccountID(ctx, other.ID)
	if err != nil || count != 1 {
		t.Fatalf("count other account: count=%d err=%v", count, err)
	}
	oldCount, err := uc.projects.CountByGitHubAccountID(ctx, account.ID)
	if err != nil || oldCount != 0 {
		t.Fatalf("count original account: count=%d err=%v", oldCount, err)
	}
}

func newTestUseCase(t *testing.T) (*UseCase, *domain.GitHubAccount, *domain.DokployServer) {
	t.Helper()
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationGitHubApp, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "demo", "https://dokploy.example.com", "server-1", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	uc := NewUseCase(
		newMemoryProjects(),
		&memoryAccounts{byID: map[uuid.UUID]*domain.GitHubAccount{account.ID: account}},
		&memoryServers{byID: map[uuid.UUID]*domain.DokployServer{server.ID: server}},
		memorySnapshots{},
	)
	uc.now = func() time.Time { return now }
	return uc, account, server
}

func assertNoSecrets(t *testing.T, project *domain.Project) {
	t.Helper()
	payload, err := json.Marshal(project)
	if err != nil {
		t.Fatal(err)
	}
	encoded := strings.ToLower(string(payload))
	for _, leaked := range []string{`"token"`, `"password"`, `"api_key"`, `"api_credential"`, `"pat"`, `"master_key"`, `"secret"`} {
		if strings.Contains(encoded, leaked) {
			t.Fatalf("project payload leaked %q: %s", leaked, payload)
		}
	}
}
