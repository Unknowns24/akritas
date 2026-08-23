package remediation_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/remediation"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// seedIncident builds the minimal GitHubAccount -> DokployServer -> Project
// -> Incident chain remediations.incident_id's foreign key requires. No
// incident repository exists yet (out of scope for this task), so the
// incidents row is inserted directly, mirroring the columns from its
// migration.
func seedIncident(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Microsecond)

	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatalf("NewGitHubAccount: %v", err)
	}
	server, err := domain.NewDokployServer(uuid.New(), "Dokploy", "https://dokploy.example.com", strings.Repeat("a", 64), domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatalf("NewDokployServer: %v", err)
	}
	if err := db.Table("github_accounts").Create(account).Error; err != nil {
		t.Fatalf("seed github_accounts: %v", err)
	}
	if err := db.Table("dokploy_servers").Create(server).Error; err != nil {
		t.Fatalf("seed dokploy_servers: %v", err)
	}

	repositoryRef, err := domain.NewGitHubRepository(account.ID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatalf("NewGitHubRepository: %v", err)
	}
	application, err := domain.NewDokployApplication(server.ID, uuid.NewString(), "api", "API", "production", domain.DokployApplicationRunning)
	if err != nil {
		t.Fatalf("NewDokployApplication: %v", err)
	}
	source, err := domain.SourceFromApplication(application)
	if err != nil {
		t.Fatalf("SourceFromApplication: %v", err)
	}
	project, err := domain.NewProject(uuid.New(), "Fixture "+uuid.NewString(), "demo", repositoryRef, source, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatalf("NewProject: %v", err)
	}
	if err := db.Table("projects").Create(project).Error; err != nil {
		t.Fatalf("seed projects: %v", err)
	}

	incidentID := uuid.New()
	err = db.Exec(`INSERT INTO incidents (
        id, key, project_id, fingerprint, severity, phase,
        first_seen_at, last_seen_at, occurrence_count, title, summary
    ) VALUES (?, ?, ?, ?, 'error', 'detected', ?, ?, 1, 'fixture incident', 'fixture')`,
		incidentID, "fixture-"+uuid.NewString(), project.ID, "fixture-fingerprint-"+uuid.NewString(), now, now,
	).Error
	if err != nil {
		t.Fatalf("seed incidents: %v", err)
	}
	return incidentID
}

func seedInvestigation(t *testing.T, db *gorm.DB, incidentID uuid.UUID) uuid.UUID {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Microsecond)
	value, err := domain.NewInvestigation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewInvestigation: %v", err)
	}
	if err := db.Table("investigations").Create(value).Error; err != nil {
		t.Fatalf("seed investigations: %v", err)
	}
	return value.ID
}

func TestRepositoryCreateAndGet(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := remediation.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	incidentID := seedIncident(t, db)
	now := time.Now().UTC().Truncate(time.Microsecond)

	value, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := repo.Create(context.Background(), value); err != nil {
		t.Fatalf("Create (planned): %v", err)
	}

	got, err := repo.Get(context.Background(), value.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != value.ID || got.IncidentID != value.IncidentID || got.Status != domain.RemediationStatusPlanned {
		t.Fatalf("unexpected round-trip: %+v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("round-tripped value fails domain validation: %v", err)
	}

	inProgress, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := inProgress.Start("akritas/remediation/"+inProgress.ID.String(), now); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := repo.Create(context.Background(), inProgress); err != nil {
		t.Fatalf("Create (in_progress): %v", err)
	}
	gotInProgress, err := repo.Get(context.Background(), inProgress.ID)
	if err != nil {
		t.Fatalf("Get (in_progress): %v", err)
	}
	if gotInProgress.Status != domain.RemediationStatusInProgress || gotInProgress.BranchName == "" {
		t.Fatalf("unexpected in_progress round-trip: %+v", gotInProgress)
	}
}

func TestRepositoryFindByInvestigationAndPullRequestRoundTrip(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := remediation.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	incidentID := seedIncident(t, db)
	investigationID := seedInvestigation(t, db, incidentID)
	now := time.Now().UTC().Truncate(time.Microsecond)

	value, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := value.AttachInvestigation(investigationID, now.Add(time.Second)); err != nil {
		t.Fatalf("AttachInvestigation: %v", err)
	}
	if err := value.Start("akritas/remediation/"+value.ID.String(), now.Add(2*time.Second)); err != nil {
		t.Fatalf("Start: %v", err)
	}
	result, err := domain.NewValidationResult(uuid.New(), value.ID, domain.ValidationTypeTest, "go test", now.Add(3*time.Second))
	if err != nil {
		t.Fatalf("NewValidationResult: %v", err)
	}
	if err := result.Start(now.Add(4 * time.Second)); err != nil {
		t.Fatalf("Start validation: %v", err)
	}
	if err := result.PassWithOutputRedacted(now.Add(5*time.Second), "ok", "safe output", false); err != nil {
		t.Fatalf("PassWithOutputRedacted: %v", err)
	}
	if err := value.AddValidationResult(*result, now.Add(6*time.Second)); err != nil {
		t.Fatalf("AddValidationResult: %v", err)
	}
	change, err := domain.NewCodeChange("main.go", domain.CodeChangeModified, "diff --git a/main.go b/main.go")
	if err != nil {
		t.Fatalf("NewCodeChange: %v", err)
	}
	if err := value.AddChange(change, now.Add(7*time.Second)); err != nil {
		t.Fatalf("AddChange: %v", err)
	}
	if err := value.MarkValidated(now.Add(8 * time.Second)); err != nil {
		t.Fatalf("MarkValidated: %v", err)
	}
	pr, err := domain.NewPullRequestReference(12, "https://github.com/Unknowns24/akritas/pull/12", "Unknowns24/akritas", value.BranchName, now.Add(9*time.Second))
	if err != nil {
		t.Fatalf("NewPullRequestReference: %v", err)
	}
	if err := value.CreatePullRequest(pr, now.Add(10*time.Second)); err != nil {
		t.Fatalf("CreatePullRequest: %v", err)
	}
	if err := repo.Create(context.Background(), value); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.FindByInvestigation(context.Background(), investigationID)
	if err != nil {
		t.Fatalf("FindByInvestigation: %v", err)
	}
	if got.ID != value.ID || got.InvestigationID != investigationID || got.PullRequestReference == nil ||
		got.PullRequestReference.Number != 12 || got.PullRequestReference.Branch != value.BranchName ||
		got.Status != domain.RemediationStatusPullRequestCreated {
		t.Fatalf("unexpected round-trip: %+v", got)
	}

	duplicate, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation duplicate: %v", err)
	}
	if err := duplicate.AttachInvestigation(investigationID, now.Add(time.Second)); err != nil {
		t.Fatalf("AttachInvestigation duplicate: %v", err)
	}
	if err := repo.Create(context.Background(), duplicate); err == nil {
		t.Fatal("expected duplicate investigation remediation to fail")
	}
}

func TestRepositoryUpdateFailureState(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := remediation.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	incidentID := seedIncident(t, db)
	now := time.Now().UTC().Truncate(time.Microsecond)

	value, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := value.Start("akritas/remediation/"+value.ID.String(), now.Add(time.Second)); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := repo.Create(context.Background(), value); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := value.Fail("La remediación falló porque una validación no pasó.", now.Add(2*time.Second)); err != nil {
		t.Fatalf("Fail: %v", err)
	}
	if err := repo.Update(context.Background(), value); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, err := repo.Get(context.Background(), value.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Status != domain.RemediationStatusFailed || got.FailureUserMessage == "" {
		t.Fatalf("unexpected updated remediation: %+v", got)
	}
}

func TestRepositoryGetNotFound(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := remediation.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = repo.Get(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrRemediationNotFound) {
		t.Fatalf("expected ErrRemediationNotFound, got %v", err)
	}
}

func TestRepositoryCreateDuplicateID(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := remediation.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	incidentID := seedIncident(t, db)
	now := time.Now().UTC().Truncate(time.Microsecond)

	value, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := repo.Create(context.Background(), value); err != nil {
		t.Fatalf("first Create: %v", err)
	}

	duplicate, err := domain.NewRemediation(value.ID, incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	err = repo.Create(context.Background(), duplicate)
	if err == nil {
		t.Fatal("expected an error for a duplicate ID")
	}
	if errors.Is(err, domain.ErrRemediationNotFound) {
		t.Fatalf("did not expect ErrRemediationNotFound for a duplicate-key failure, got %v", err)
	}
}
