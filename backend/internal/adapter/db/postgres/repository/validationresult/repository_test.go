package validationresult_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	remediationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/remediation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/validationresult"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// seedRemediation builds the GitHubAccount -> DokployServer -> Project ->
// Incident -> Remediation chain validation_results.remediation_id's
// foreign key requires.
func seedRemediation(t *testing.T, db *gorm.DB) uuid.UUID {
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

	remediationRepo, err := remediationrepo.New(db)
	if err != nil {
		t.Fatalf("remediation.New: %v", err)
	}
	value, err := domain.NewRemediation(uuid.New(), incidentID, now)
	if err != nil {
		t.Fatalf("NewRemediation: %v", err)
	}
	if err := value.Start("akritas/remediation/"+value.ID.String(), now); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := remediationRepo.Create(context.Background(), value); err != nil {
		t.Fatalf("seed remediations: %v", err)
	}
	return value.ID
}

func newResult(t *testing.T, remediationID uuid.UUID, validationType domain.ValidationType, status domain.ValidationStatus, now time.Time) domain.ValidationResult {
	t.Helper()
	value, err := domain.NewValidationResult(uuid.New(), remediationID, validationType, string(validationType)+" step", now)
	if err != nil {
		t.Fatalf("NewValidationResult: %v", err)
	}
	if status == domain.ValidationStatusPending {
		return *value
	}
	if err := value.Start(now.Add(time.Second)); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if status == domain.ValidationStatusRunning {
		return *value
	}
	if status == domain.ValidationStatusPassed {
		if err := value.Pass(now.Add(2*time.Second), "ok", "output"); err != nil {
			t.Fatalf("Pass: %v", err)
		}
		return *value
	}
	if err := value.Fail(now.Add(2*time.Second), "failed", "output"); err != nil {
		t.Fatalf("Fail: %v", err)
	}
	return *value
}

func TestRepositoryCreateAndListByRemediation(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := validationresult.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	remediationID := seedRemediation(t, db)
	now := time.Now().UTC().Truncate(time.Microsecond)

	values := []domain.ValidationResult{
		newResult(t, remediationID, domain.ValidationTypeBuild, domain.ValidationStatusPassed, now),
		newResult(t, remediationID, domain.ValidationTypeStaticAnalysis, domain.ValidationStatusFailed, now.Add(10*time.Second)),
		newResult(t, remediationID, domain.ValidationTypeTest, domain.ValidationStatusPending, now.Add(20*time.Second)),
	}
	for i := range values {
		if err := repo.Create(context.Background(), &values[i]); err != nil {
			t.Fatalf("Create[%d]: %v", i, err)
		}
	}

	got, err := repo.ListByRemediation(context.Background(), remediationID)
	if err != nil {
		t.Fatalf("ListByRemediation: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	for i, value := range got {
		if err := value.Validate(); err != nil {
			t.Fatalf("result %d fails domain validation: %v (%+v)", i, err, value)
		}
		if value.ID != values[i].ID {
			t.Fatalf("expected stable creation order, got %+v at index %d, want ID %s", value, i, values[i].ID)
		}
	}
}

func TestRepositoryListByRemediationEmpty(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := validationresult.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	got, err := repo.ListByRemediation(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("ListByRemediation: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no results, got %d", len(got))
	}
}

func TestRepositoryCreateForeignKeyViolation(t *testing.T) {
	db := dbtest.Connect(t)
	repo, err := validationresult.New(db)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Now().UTC().Truncate(time.Microsecond)
	value := newResult(t, uuid.New(), domain.ValidationTypeBuild, domain.ValidationStatusPending, now)

	err = repo.Create(context.Background(), &value)
	if err == nil {
		t.Fatal("expected a foreign key violation error")
	}
	if errors.Is(err, domain.ErrRemediationNotFound) {
		t.Fatalf("did not expect ErrRemediationNotFound for a FK violation, got %v", err)
	}
}
