package project

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func TestProjectUseCaseCreateActivateDisableAndDelete(t *testing.T) {
	t.Parallel()
	fixture := newFixture(t)
	created, err := fixture.uc.Create(context.Background(), fixture.command("Akritas", "app-1"))
	if err != nil {
		t.Fatal(err)
	}
	if created.Project.MonitoringStatus != domain.MonitoringStatusDisabled || len(created.BuiltInDetectionRules) != 7 {
		t.Fatalf("create result: %+v", created)
	}

	enabled, _ := domain.NewMonitoringConfiguration(true, []string{"panic"}, []string{"healthcheck"}, 15*time.Minute, 5, 6)
	saved, err := fixture.uc.PutMonitoring(context.Background(), created.Project.ID, enabled)
	if err != nil {
		t.Fatal(err)
	}
	if !saved.Enabled || fixture.store.value.MonitoringStatus != domain.MonitoringStatusStarting || fixture.github.calls != 2 || fixture.dokploy.calls != 2 {
		t.Fatalf("activation did not revalidate: saved=%+v", saved)
	}
	if err := fixture.uc.Delete(context.Background(), created.Project.ID); !errors.Is(err, domain.ErrProjectMustBeDisabled) {
		t.Fatalf("active delete = %v", err)
	}

	disabled := enabled.Clone()
	disabled.Enabled = false
	if _, err := fixture.uc.PutMonitoring(context.Background(), created.Project.ID, disabled); err != nil {
		t.Fatal(err)
	}
	if err := fixture.uc.Delete(context.Background(), created.Project.ID); err != nil {
		t.Fatal(err)
	}
	if fixture.store.value != nil {
		t.Fatal("project survived delete")
	}
}

func TestProjectUseCaseConflictsAndProviderFailureDoNotMutate(t *testing.T) {
	t.Parallel()
	fixture := newFixture(t)
	created, err := fixture.uc.Create(context.Background(), fixture.command("Akritas", "app-1"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.uc.Create(context.Background(), fixture.command("akritas", "app-2")); !errors.Is(err, domain.ErrProjectNameConflict) {
		t.Fatalf("duplicate name = %v", err)
	}

	before := *fixture.store.value
	fixture.github.err = domain.ErrIntegrationUnavailable
	enabled := domain.DefaultMonitoringConfiguration()
	enabled.Enabled = true
	if _, err := fixture.uc.PutMonitoring(context.Background(), created.Project.ID, enabled); !errors.Is(err, domain.ErrIntegrationUnavailable) {
		t.Fatalf("provider failure = %v", err)
	}
	if fixture.store.updates != 0 || fixture.store.value.MonitoringStatus != before.MonitoringStatus {
		t.Fatal("failed activation mutated persistence")
	}
}

func TestProjectUseCaseRejectsExclusiveApplicationAndDefaultBranchMismatch(t *testing.T) {
	t.Parallel()
	fixture := newFixture(t)
	if _, err := fixture.uc.Create(context.Background(), fixture.command("Akritas", "app-1")); err != nil {
		t.Fatal(err)
	}
	if _, err := fixture.uc.Create(context.Background(), fixture.command("Another", "app-1")); !errors.Is(err, domain.ErrProjectDokploySourceConflict) {
		t.Fatalf("exclusive application = %v", err)
	}

	other := newFixture(t)
	other.github.branch = "develop"
	if _, err := other.uc.Create(context.Background(), other.command("Akritas", "app-1")); !errors.Is(err, domain.ErrProjectDefaultBranchMismatch) {
		t.Fatalf("default branch mismatch = %v", err)
	}
	if other.store.value != nil {
		t.Fatal("branch mismatch persisted a Project")
	}
}

func TestProjectUseCaseKeepsActiveAssociationUpdatesAtomic(t *testing.T) {
	t.Parallel()
	fixture := newFixture(t)
	created, err := fixture.uc.Create(context.Background(), fixture.command("Akritas", "app-1"))
	if err != nil {
		t.Fatal(err)
	}
	enabled := domain.DefaultMonitoringConfiguration()
	enabled.Enabled = true
	if _, err := fixture.uc.PutMonitoring(context.Background(), created.Project.ID, enabled); err != nil {
		t.Fatal(err)
	}
	updatesBefore := fixture.store.updates
	newName, repositoryIdentifier := "Renamed", "84"
	_, err = fixture.uc.Update(context.Background(), portsin.UpdateProjectCommand{ID: created.Project.ID, Name: &newName, RepositoryIdentifier: &repositoryIdentifier})
	if !errors.Is(err, domain.ErrProjectMustBeDisabled) {
		t.Fatalf("active association update = %v", err)
	}
	if fixture.store.updates != updatesBefore || fixture.store.value.Name != "Akritas" || fixture.store.value.GitHubRepository.RepositoryIdentifier != "42" {
		t.Fatalf("rejected update was not atomic: %+v", fixture.store.value)
	}
}

func TestProjectUseCaseReadsAndUpdatesDisabledProject(t *testing.T) {
	t.Parallel()
	fixture := newFixture(t)
	created, err := fixture.uc.Create(context.Background(), fixture.command("Akritas", "app-1"))
	if err != nil {
		t.Fatal(err)
	}
	name, description := "Renamed", "updated"
	selector := domain.DokploySourceSelector{Type: domain.DokploySourceApplication, DokployServerID: fixture.server.ID, ResourceIdentifier: "app-2"}
	updated, err := fixture.uc.Update(context.Background(), portsin.UpdateProjectCommand{ID: created.Project.ID, Name: &name, Description: &description, DokploySource: &selector})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Project.Name != name || updated.Project.Description != description || updated.Project.DokploySource.ResourceIdentifier != "app-2" {
		t.Fatalf("updated Project = %+v", updated.Project)
	}
	got, err := fixture.uc.Get(context.Background(), created.Project.ID)
	if err != nil || got.Project.Name != name {
		t.Fatalf("Get() = %+v, %v", got, err)
	}
	monitoring, err := fixture.uc.GetMonitoring(context.Background(), created.Project.ID)
	if err != nil || monitoring.Enabled {
		t.Fatalf("GetMonitoring() = %+v, %v", monitoring, err)
	}
	page, err := fixture.uc.List(context.Background(), paging.Params{Limit: 25})
	if err != nil || page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("List() = %+v, %v", page, err)
	}

	updatesBefore := fixture.store.updates
	repositoryIdentifier := "84"
	fixture.github.err = domain.ErrIntegrationUnavailable
	if _, err := fixture.uc.Update(context.Background(), portsin.UpdateProjectCommand{ID: created.Project.ID, RepositoryIdentifier: &repositoryIdentifier}); !errors.Is(err, domain.ErrIntegrationUnavailable) {
		t.Fatalf("provider update failure = %v", err)
	}
	if fixture.store.updates != updatesBefore || fixture.store.value.GitHubRepository.RepositoryIdentifier != "42" {
		t.Fatal("failed provider resolution mutated persistence")
	}
}

func TestProjectUseCaseCreatesComposeServiceAndRejectsMissingService(t *testing.T) {
	fixture := newFixture(t)
	command := fixture.command("Compose API", "unused")
	command.DokploySource = domain.DokploySourceSelector{Type: domain.DokploySourceComposeService, DokployServerID: fixture.server.ID, ResourceIdentifier: "compose-1", ServiceName: "api"}
	created, err := fixture.uc.Create(context.Background(), command)
	if err != nil {
		t.Fatal(err)
	}
	if created.Project.DokploySource.Type != domain.DokploySourceComposeService || created.Project.DokploySource.ServiceName != "api" || created.Project.DokploySource.RuntimeType != domain.DokployRuntimeCompose {
		t.Fatalf("source = %+v", created.Project.DokploySource)
	}

	missing := newFixture(t)
	missingCommand := missing.command("Missing", "unused")
	missingCommand.DokploySource = domain.DokploySourceSelector{Type: domain.DokploySourceComposeService, DokployServerID: missing.server.ID, ResourceIdentifier: "compose-1", ServiceName: "worker"}
	if _, err := missing.uc.Create(context.Background(), missingCommand); !errors.Is(err, domain.ErrProjectDokploySourceNotFound) {
		t.Fatalf("missing service = %v", err)
	}
	if missing.store.value != nil {
		t.Fatal("missing service persisted Project")
	}
}

type fixture struct {
	uc      *UseCase
	store   *memoryStore
	github  *repositoryResolver
	dokploy *applicationResolver
	account *domain.GitHubAccount
	server  *domain.DokployServer
}

func newFixture(t *testing.T) fixture {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "Dokploy", "https://dokploy.example.com", strings.Repeat("a", 64), domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	store := &memoryStore{}
	github := &repositoryResolver{}
	dokploy := &applicationResolver{}
	input := New(store, accountReader{account}, serverReader{server}, github, dokploy, func() uuid.UUID { return uuid.MustParse("00000000-0000-0000-0000-000000000123") }, func() time.Time { return now })
	return fixture{uc: input.(*UseCase), store: store, github: github, dokploy: dokploy, account: account, server: server}
}

func (f fixture) command(name, app string) portsin.CreateProjectCommand {
	return portsin.CreateProjectCommand{Name: name, GitHubAccountID: f.account.ID, RepositoryIdentifier: "42", DefaultBranch: "main", DokploySource: domain.DokploySourceSelector{Type: domain.DokploySourceApplication, DokployServerID: f.server.ID, ResourceIdentifier: app}, MonitoringConfiguration: domain.DefaultMonitoringConfiguration()}
}

type accountReader struct{ value *domain.GitHubAccount }

func (r accountReader) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return r.value, nil
}

type serverReader struct{ value *domain.DokployServer }

func (r serverReader) Get(context.Context, uuid.UUID) (*domain.DokployServer, error) {
	return r.value, nil
}

type repositoryResolver struct {
	calls  int
	err    error
	branch string
}

func (r *repositoryResolver) GetRepository(_ context.Context, account domain.GitHubAccount, identifier string) (domain.GitHubRepository, error) {
	r.calls++
	if r.err != nil {
		return domain.GitHubRepository{}, r.err
	}
	branch := r.branch
	if branch == "" {
		branch = "main"
	}
	return domain.NewGitHubRepository(account.ID, identifier, account.AccountIdentifier, "akritas", branch, true, "https://github.com/Unknowns24/akritas")
}

type applicationResolver struct {
	calls int
	err   error
}

func (r *applicationResolver) GetApplication(_ context.Context, server domain.DokployServer, identifier string) (domain.DokployApplication, error) {
	r.calls++
	if r.err != nil {
		return domain.DokployApplication{}, r.err
	}
	return domain.NewDokployApplication(server.ID, identifier, "instance-"+identifier, identifier, "production", domain.DokployApplicationRunning)
}

func (r *applicationResolver) GetCompose(_ context.Context, server domain.DokployServer, identifier string) (domain.DokployCompose, error) {
	return domain.NewDokployCompose(server.ID, identifier, "stack-"+identifier, identifier, "env", "production", domain.DokploySourceRunning, domain.DokployRuntimeCompose, "")
}

func (r *applicationResolver) ListComposeServices(_ context.Context, server domain.DokployServer, identifier string, _ bool) ([]domain.DokployComposeService, error) {
	service, err := domain.NewDokployComposeService(server.ID, identifier, "api")
	return []domain.DokployComposeService{service}, err
}

type memoryStore struct {
	value   *domain.Project
	updates int
}

func (s *memoryStore) Create(_ context.Context, value *domain.Project) error {
	copy := *value
	s.value = &copy
	return nil
}
func (s *memoryStore) Get(_ context.Context, id uuid.UUID) (*domain.Project, error) {
	if s.value == nil || s.value.ID != id {
		return nil, domain.ErrProjectNotFound
	}
	copy := *s.value
	return &copy, nil
}
func (s *memoryStore) FindByNormalizedName(_ context.Context, name string) (*domain.Project, error) {
	if s.value != nil && strings.EqualFold(s.value.Name, strings.TrimSpace(name)) {
		copy := *s.value
		return &copy, nil
	}
	return nil, domain.ErrProjectNotFound
}
func (s *memoryStore) FindByDokploySource(_ context.Context, selector domain.DokploySourceSelector) (*domain.Project, error) {
	if s.value != nil && s.value.DokploySource.IdentityKey() == selector.IdentityKey() {
		copy := *s.value
		return &copy, nil
	}
	return nil, domain.ErrProjectNotFound
}
func (s *memoryStore) List(context.Context, paging.Params) (paging.Slice[domain.Project], error) {
	if s.value == nil {
		return paging.Slice[domain.Project]{}, nil
	}
	return paging.Slice[domain.Project]{Items: []domain.Project{*s.value}, Total: 1}, nil
}
func (s *memoryStore) Update(_ context.Context, value *domain.Project, expected time.Time) error {
	if s.value == nil || s.value.UpdatedAt != expected {
		return domain.ErrProjectConcurrentModification
	}
	copy := *value
	s.value = &copy
	s.updates++
	return nil
}
func (s *memoryStore) Delete(_ context.Context, id uuid.UUID) error {
	if s.value == nil || s.value.ID != id {
		return domain.ErrProjectNotFound
	}
	s.value = nil
	return nil
}
