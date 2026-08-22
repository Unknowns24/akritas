package postgres_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestProjectRepositoryRoundTripAndConstraints(t *testing.T) {
	db := dbtest.OpenMigrated(t)
	ctx := context.Background()
	now := time.Date(2026, 8, 22, 16, 0, 0, 0, time.UTC)
	accounts := githubaccount.NewRepository(db)
	servers := dokployserver.NewRepository(db)
	projects := project.NewRepository(db)

	account, server := seedLookups(t, ctx, accounts, servers, now)
	created := newPersistedProject(t, account, server, "sentinel-api", "app-1", now)
	if err := projects.Create(ctx, created); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := projects.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "sentinel-api" || got.GitHubRepository.FullName != "Unknowns24/akritas" {
		t.Fatalf("round trip mismatch: %+v", got)
	}
	if got.GitHubRepository.GitHubAccountID != account.ID ||
		got.GitHubRepository.RepositoryIdentifier != "Unknowns24/akritas" ||
		got.GitHubRepository.Owner != "Unknowns24" ||
		got.GitHubRepository.Name != "akritas" ||
		got.GitHubRepository.DefaultBranch != "main" ||
		got.GitHubRepository.Private ||
		got.GitHubRepository.HTMLURL != "https://github.com/Unknowns24/akritas" {
		t.Fatalf("github snapshot round trip mismatch: %+v", got.GitHubRepository)
	}
	if got.DokployApplication.DokployServerID != server.ID ||
		got.DokployApplication.ApplicationIdentifier != "app-1" ||
		got.DokployApplication.InstanceIdentifier != "app-1" ||
		got.DokployApplication.DisplayName != "app-1" ||
		got.DokployApplication.Environment != "production" ||
		got.DokployApplication.Status != domain.DokployApplicationUnknown {
		t.Fatalf("dokploy snapshot round trip mismatch: %+v", got.DokployApplication)
	}
	byApp, err := projects.GetByDokployApplication(ctx, server.ID, "app-1")
	if err != nil || byApp.ID != created.ID {
		t.Fatalf("get by dokploy application: %+v err=%v", byApp, err)
	}

	count, err := projects.CountByGitHubAccountID(ctx, account.ID)
	if err != nil || count != 1 {
		t.Fatalf("count by account: count=%d err=%v", count, err)
	}
	emptyCount, err := projects.CountByGitHubAccountID(ctx, uuid.New())
	if err != nil || emptyCount != 0 {
		t.Fatalf("count missing account: count=%d err=%v", emptyCount, err)
	}

	second := newPersistedProject(t, account, server, "sentinel-worker", "app-2", now)
	if err := projects.Create(ctx, second); err != nil {
		t.Fatalf("create second: %v", err)
	}
	count, err = projects.CountByGitHubAccountID(ctx, account.ID)
	if err != nil || count != 2 {
		t.Fatalf("count after second project: count=%d err=%v", count, err)
	}

	orphan := newPersistedProject(t, account, server, "orphan", "app-orphan", now)
	orphan.GitHubRepository.GitHubAccountID = uuid.New()
	if err := projects.Create(ctx, orphan); err == nil {
		t.Fatal("expected FK violation for unknown github_account_id")
	}

	serverCount, err := projects.CountByDokployServerID(ctx, server.ID)
	if err != nil || serverCount != 2 {
		t.Fatalf("count by dokploy server: count=%d err=%v", serverCount, err)
	}
	emptyServerCount, err := projects.CountByDokployServerID(ctx, uuid.New())
	if err != nil || emptyServerCount != 0 {
		t.Fatalf("count missing server: count=%d err=%v", emptyServerCount, err)
	}

	orphanServer := newPersistedProject(t, account, server, "orphan-server", "app-orphan-server", now)
	orphanServer.DokployApplication.DokployServerID = uuid.New()
	if err := projects.Create(ctx, orphanServer); err == nil {
		t.Fatal("expected FK violation for unknown dokploy_server_id")
	}

	listed, total, err := projects.List(ctx, paging.ListQuery{Limit: 10, NameLike: "sentinel", Offset: 0})
	if err != nil || total != 2 || len(listed) != 2 {
		t.Fatalf("list: total=%d n=%d err=%v", total, len(listed), err)
	}

	created.Name = "sentinel-core"
	if err := created.Rename("sentinel-core", "updated", now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := projects.Update(ctx, created); err != nil {
		t.Fatal(err)
	}
	byName, err := projects.GetByNormalizedName(ctx, "SENTINEL-CORE")
	if err != nil || byName.ID != created.ID {
		t.Fatalf("normalized name lookup failed: %v", err)
	}

	duplicateName := newPersistedProject(t, account, server, "sentinel-core", "app-2", now)
	if err := projects.Create(ctx, duplicateName); err == nil {
		t.Fatal("expected unique name constraint")
	}
	duplicateApp := newPersistedProject(t, account, server, "other", "app-1", now)
	if err := projects.Create(ctx, duplicateApp); err == nil {
		t.Fatal("expected unique dokploy application constraint")
	}

	assertSchemaHasNoSecrets(t, db)
	assertNoGitHubRepositoriesTable(t, db)
	assertNoDokployApplicationsTable(t, db)
	if db.Migrator().HasTable("monitoring_configurations") {
		t.Fatal("monitoring_configurations table must not exist; configuration is embedded on projects")
	}
	if _, err := projects.GetByID(ctx, uuid.New()); !errors.Is(err, apperr.ErrProjectNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestProjectRepositoryPersistsMonitoringConfiguration(t *testing.T) {
	db := dbtest.OpenMigrated(t)
	ctx := context.Background()
	now := time.Date(2026, 8, 22, 16, 0, 0, 0, time.UTC)
	accounts := githubaccount.NewRepository(db)
	servers := dokployserver.NewRepository(db)
	projects := project.NewRepository(db)
	account, server := seedLookups(t, ctx, accounts, servers, now)

	created := newPersistedProject(t, account, server, "sentinel-api", "app-1", now)
	custom, err := domain.NewMonitoringConfiguration(
		true, []string{`database .* unavailable`}, []string{`expected healthcheck failure`},
		15*time.Minute, 8, 12,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := created.ReplaceMonitoringConfiguration(custom, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := projects.Create(ctx, created); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := projects.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	assertMonitoringEquals(t, got.MonitoringConfiguration, custom)

	cleared := domain.DefaultMonitoringConfiguration()
	if err := got.ReplaceMonitoringConfiguration(cleared, now.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := projects.Update(ctx, got); err != nil {
		t.Fatal(err)
	}
	updated, err := projects.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	assertMonitoringEquals(t, updated.MonitoringConfiguration, cleared)
}

func assertMonitoringEquals(t *testing.T, got, want domain.MonitoringConfiguration) {
	t.Helper()
	if got.Enabled != want.Enabled || got.GroupingWindow != want.GroupingWindow ||
		got.ContextBefore != want.ContextBefore || got.ContextAfter != want.ContextAfter {
		t.Fatalf("monitoring mismatch: got=%+v want=%+v", got, want)
	}
	if len(got.ErrorPatterns) != len(want.ErrorPatterns) || len(got.IgnoredPatterns) != len(want.IgnoredPatterns) {
		t.Fatalf("pattern lengths: got=%+v want=%+v", got, want)
	}
	for i := range want.ErrorPatterns {
		if got.ErrorPatterns[i] != want.ErrorPatterns[i] {
			t.Fatalf("error pattern %d: %q vs %q", i, got.ErrorPatterns[i], want.ErrorPatterns[i])
		}
	}
	for i := range want.IgnoredPatterns {
		if got.IgnoredPatterns[i] != want.IgnoredPatterns[i] {
			t.Fatalf("ignored pattern %d: %q vs %q", i, got.IgnoredPatterns[i], want.IgnoredPatterns[i])
		}
	}
}

func seedLookups(
	t *testing.T,
	ctx context.Context,
	accounts *githubaccount.Repository,
	servers *dokployserver.Repository,
	now time.Time,
) (*domain.GitHubAccount, *domain.DokployServer) {
	t.Helper()
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationGitHubApp, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "demo", "https://dokploy.example.com", "server-1", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := accounts.Create(ctx, account); err != nil {
		t.Fatal(err)
	}
	if err := servers.Create(ctx, server); err != nil {
		t.Fatal(err)
	}
	return account, server
}

func newPersistedProject(t *testing.T, account *domain.GitHubAccount, server *domain.DokployServer, name, appID string, now time.Time) *domain.Project {
	t.Helper()
	repository, err := domain.NewGitHubRepository(account.ID, "Unknowns24/akritas", "Unknowns24", "akritas", "main", false, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	application, err := domain.NewDokployApplication(server.ID, appID, appID, appID, "production", domain.DokployApplicationUnknown)
	if err != nil {
		t.Fatal(err)
	}
	project, err := domain.NewProject(uuid.New(), name, "", repository, application, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	return project
}

func assertNoGitHubRepositoriesTable(t *testing.T, db *gorm.DB) {
	t.Helper()
	if db.Migrator().HasTable("github_repositories") {
		t.Fatal("github_repositories table must not exist; repository is an embedded value object on projects")
	}
}

func assertNoDokployApplicationsTable(t *testing.T, db *gorm.DB) {
	t.Helper()
	if db.Migrator().HasTable("dokploy_applications") {
		t.Fatal("dokploy_applications table must not exist; application is an embedded value object on projects")
	}
}

func assertSchemaHasNoSecrets(t *testing.T, db *gorm.DB) {
	t.Helper()
	rows, err := db.Raw(`
		SELECT table_name || '.' || column_name
		FROM information_schema.columns
		WHERE table_schema = current_schema()
		UNION ALL
		SELECT indexname || ':' || indexdef
		FROM pg_indexes
		WHERE schemaname = current_schema()
	`).Rows()
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	denied := []string{"token", "password", "secret", "api_key", "api_credential", "master_key"}
	for rows.Next() {
		var ddl string
		if err := rows.Scan(&ddl); err != nil {
			t.Fatal(err)
		}
		lower := strings.ToLower(ddl)
		for _, word := range denied {
			if strings.Contains(lower, word) {
				t.Fatalf("schema leaked %q in %s", word, ddl)
			}
		}
	}
}
