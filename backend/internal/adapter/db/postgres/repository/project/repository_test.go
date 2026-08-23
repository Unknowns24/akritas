package project

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestRepositoryPersistsSnapshotsEnforcesUsageAndOptimisticUpdate(t *testing.T) {
	db := dbtest.Connect(t)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "Dokploy", "https://dokploy.example.com", strings.Repeat("a", 64), domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Table("github_accounts").Create(account).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table("dokploy_servers").Create(server).Error; err != nil {
		t.Fatal(err)
	}
	repository, _ := domain.NewGitHubRepository(account.ID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	application, _ := domain.NewDokployApplication(server.ID, "app-1", "api", "API", "production", domain.DokployApplicationRunning)
	source, _ := domain.SourceFromApplication(application)
	value, err := domain.NewProject(uuid.New(), "Akritas", "demo", repository, source, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	store, _ := New(db)
	if err := store.Create(ctx, value); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(ctx, value.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.GitHubRepository.RepositoryIdentifier != "42" || got.DokploySource.Environment != "production" || got.MonitoringConfiguration.GroupingWindow != 30*time.Minute {
		t.Fatalf("round trip = %+v", got)
	}
	if used, err := store.GitHubAccountInUse(ctx, account.ID); err != nil || !used {
		t.Fatalf("github usage = %v, %v", used, err)
	}
	if used, err := store.DokployServerInUse(ctx, server.ID); err != nil || !used {
		t.Fatalf("dokploy usage = %v, %v", used, err)
	}

	expected := got.UpdatedAt
	if err := got.Rename("Renamed", "demo", now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := store.Update(ctx, got, expected); err != nil {
		t.Fatal(err)
	}
	if err := store.Update(ctx, got, expected); !errors.Is(err, domain.ErrProjectConcurrentModification) {
		t.Fatalf("stale update = %v", err)
	}

	duplicateName := *got
	duplicateName.ID = uuid.New()
	duplicateName.Name = "RENAMED"
	duplicateName.DokploySource.ResourceIdentifier = "app-2"
	if err := store.Create(ctx, &duplicateName); !errors.Is(err, domain.ErrProjectConcurrentModification) {
		t.Fatalf("case-insensitive name constraint = %v", err)
	}
	duplicateApplication := *got
	duplicateApplication.ID = uuid.New()
	duplicateApplication.Name = "Another"
	if err := store.Create(ctx, &duplicateApplication); !errors.Is(err, domain.ErrProjectConcurrentModification) {
		t.Fatalf("exclusive application constraint = %v", err)
	}
	composeAPI, _ := domain.NewDokploySource(server.ID, domain.DokploySourceComposeService, "compose-1", "api", "stack-prod", "Platform / api", "production", domain.DokploySourceRunning, domain.DokployRuntimeCompose, "")
	composeWorker, _ := domain.NewDokploySource(server.ID, domain.DokploySourceComposeService, "compose-1", "worker", "stack-prod", "Platform / worker", "production", domain.DokploySourceRunning, domain.DokployRuntimeCompose, "")
	apiProject, _ := domain.NewProject(uuid.New(), "Compose API", "", repository, composeAPI, domain.DefaultMonitoringConfiguration(), now)
	workerProject, _ := domain.NewProject(uuid.New(), "Compose Worker", "", repository, composeWorker, domain.DefaultMonitoringConfiguration(), now)
	if err := store.Create(ctx, apiProject); err != nil {
		t.Fatalf("create compose API = %v", err)
	}
	if err := store.Create(ctx, workerProject); err != nil {
		t.Fatalf("different service in same Compose = %v", err)
	}
	duplicateService := *apiProject
	duplicateService.ID = uuid.New()
	duplicateService.Name = "Compose API duplicate"
	if err := store.Create(ctx, &duplicateService); !errors.Is(err, domain.ErrProjectConcurrentModification) {
		t.Fatalf("exclusive compose service constraint = %v", err)
	}
	missingAccount := *got
	missingAccount.ID = uuid.New()
	missingAccount.Name = "Missing account"
	missingAccount.GitHubRepository.GitHubAccountID = uuid.New()
	missingAccount.DokploySource.ResourceIdentifier = "app-3"
	if err := store.Create(ctx, &missingAccount); err == nil {
		t.Fatal("foreign key constraint accepted a missing GitHub account")
	}

	page, err := store.List(ctx, paging.Params{
		Limit:   25,
		Filters: map[string]string{"name_like": "Ren", "monitoring_status_in": "disabled"},
		Sort:    []ukerpagination.SortExpression{{Field: "name", Direction: ukerpagination.DirectionAsc}},
	})
	if err != nil || page.Total != 1 || len(page.Items) != 1 || page.Items[0].Name != "Renamed" {
		t.Fatalf("filtered list = %+v, %v", page, err)
	}
}
