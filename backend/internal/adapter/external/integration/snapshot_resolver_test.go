package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestSnapshotResolverParsesIdentifiers(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 17, 0, 0, 0, time.UTC)
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationGitHubApp, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "demo", "https://dokploy.example.com", "server-1", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	resolver := NewSnapshotResolver()

	repository, err := resolver.ResolveGitHubRepository(account, "Unknowns24/akritas", "main")
	if err != nil {
		t.Fatal(err)
	}
	if repository.Owner != "Unknowns24" || repository.Name != "akritas" || repository.FullName != "Unknowns24/akritas" {
		t.Fatalf("unexpected repository: %+v", repository)
	}
	if repository.HTMLURL != "https://github.com/Unknowns24/akritas" {
		t.Fatalf("unexpected html url: %s", repository.HTMLURL)
	}

	simple, err := resolver.ResolveGitHubRepository(account, "sentinel", "develop")
	if err != nil {
		t.Fatal(err)
	}
	if simple.Owner != "Unknowns24" || simple.Name != "sentinel" {
		t.Fatalf("simple identifier should use account owner: %+v", simple)
	}

	if _, err := resolver.ResolveGitHubRepository(account, "a/b/c", "main"); !errors.Is(err, apperr.ErrRepositoryNotResolvable) {
		t.Fatalf("expected unresolvable repo, got %v", err)
	}
	if _, err := resolver.ResolveGitHubRepository(account, "", "main"); !errors.Is(err, apperr.ErrRepositoryNotResolvable) {
		t.Fatalf("expected empty repo error, got %v", err)
	}

	application, err := resolver.ResolveDokployApplication(server, "app-1")
	if err != nil {
		t.Fatal(err)
	}
	if application.InstanceIdentifier != "app-1" || application.DisplayName != "app-1" || application.Status != domain.DokployApplicationUnknown {
		t.Fatalf("unexpected application: %+v", application)
	}
	if _, err := resolver.ResolveDokployApplication(server, "  "); !errors.Is(err, apperr.ErrApplicationNotResolvable) {
		t.Fatalf("expected unresolvable application, got %v", err)
	}
}
