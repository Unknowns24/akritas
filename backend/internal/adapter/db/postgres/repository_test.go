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

	listed, total, err := projects.List(ctx, paging.ListQuery{Limit: 10, NameLike: "sentinel", Offset: 0})
	if err != nil || total != 1 || len(listed) != 1 {
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
	if _, err := projects.GetByID(ctx, uuid.New()); !errors.Is(err, apperr.ErrProjectNotFound) {
		t.Fatalf("expected not found, got %v", err)
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
