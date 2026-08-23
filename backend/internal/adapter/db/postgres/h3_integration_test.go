//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	dbadapter "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations/schema"
	evidencerepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/evidence"
	githubissuereferencerepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubissuereference"
	incidentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/incident"
	investigationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/investigation"
	monitoringrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/monitoring"
	operationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/operation"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/Unknowns24/akritas/backend/internal/service/evidenceassembly"
	"github.com/Unknowns24/akritas/backend/internal/service/issuecontent"
	evidenceusecase "github.com/Unknowns24/akritas/backend/internal/usecase/evidence"
	investigationusecase "github.com/Unknowns24/akritas/backend/internal/usecase/investigation"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type recordingDispatcher struct{ investigationID, operationID uuid.UUID }

func (d *recordingDispatcher) Dispatch(investigationID, operationID uuid.UUID) {
	d.investigationID, d.operationID = investigationID, operationID
}

type failingCreateOperationStore struct{ portsout.OperationStore }

func (f failingCreateOperationStore) Create(context.Context, *domain.Operation) error {
	return errors.New("simulated operation create failure")
}

type fixedAccountReader struct{ account *domain.GitHubAccount }

func (r fixedAccountReader) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return r.account, nil
}

type fixedIssuePublisher struct {
	number    int
	url       string
	createdAt time.Time
}

func (p fixedIssuePublisher) PublishIssue(context.Context, domain.GitHubAccount, domain.GitHubRepository, portsout.IssueContent) (portsout.PublishedIssue, error) {
	return portsout.PublishedIssue{Number: p.number, URL: p.url, CreatedAt: p.createdAt}, nil
}

type evidenceAwareRunner struct {
	context portsout.InvestigationRunContext
	now     time.Time
}

func (r *evidenceAwareRunner) Run(_ context.Context, runContext portsout.InvestigationRunContext) (portsout.InvestigationRunResult, error) {
	r.context = runContext
	var cited uuid.UUID
	for _, evidence := range runContext.Evidence {
		if evidence.Type == domain.EvidenceLogExcerpt && strings.Contains(evidence.Content, "database connection refused") {
			cited = evidence.ID
			break
		}
	}
	if cited == uuid.Nil {
		return portsout.InvestigationRunResult{}, errors.New("QVAC did not receive H2 log Evidence")
	}
	toolEvidence, err := domain.NewEvidence(uuid.New(), runContext.Investigation.ID, domain.EvidenceCodeLocation,
		"Repository search located database connection code.", `{"path":"internal/db/connect.go"}`, r.now)
	if err != nil {
		return portsout.InvestigationRunResult{}, err
	}
	toolEvidence.FilePath = "internal/db/connect.go"
	if err := toolEvidence.Validate(); err != nil {
		return portsout.InvestigationRunResult{}, err
	}
	return portsout.InvestigationRunResult{
		Summary: "The database endpoint refused the worker connection.", RootCause: "Database endpoint unavailable",
		RootCauseStatus: domain.RootCauseIdentified, ResolutionStatus: domain.ResolutionFixable, Confidence: 0.93,
		Hypotheses: []string{"database service was unavailable"}, RelevantFiles: []string{"internal/db/connect.go"},
		RecommendedActions: []string{"verify database availability and connection configuration"},
		EvidenceIDs:        []uuid.UUID{cited, toolEvidence.ID}, DiscoveredEvidence: []domain.Evidence{*toolEvidence},
	}, nil
}

func TestH2IncidentToPersistedH3ResultAgainstPostgreSQL(t *testing.T) {
	db := dbtest.ConnectContainer(t)
	ctx := context.Background()
	base := time.Now().UTC().Truncate(time.Second)
	account, project := insertH3Project(t, db, base)

	monitoring, err := monitoringrepo.New(db)
	if err != nil {
		t.Fatal(err)
	}
	incident, err := domain.NewIncident(uuid.New(), "AKR-700", project.ID,
		"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", domain.SeverityError,
		"database connection refused", base)
	if err != nil {
		t.Fatal(err)
	}
	if err := monitoring.CreateIncident(ctx, incident); err != nil {
		t.Fatal(err)
	}
	before, _ := domain.NewSanitizedLogRecord(base.Add(-time.Second), domain.LogStreamStderr, "opening database connection")
	after, _ := domain.NewSanitizedLogRecord(base.Add(time.Second), domain.LogStreamStderr, "retry scheduled in 5 seconds")
	event, err := domain.NewLogEvent(uuid.New(), project.ID, base, domain.SeverityError, "database connection refused",
		incident.Fingerprint, []string{string(domain.DetectionRuleErrorLevel)}, []domain.SanitizedLogRecord{before}, []domain.SanitizedLogRecord{after})
	if err != nil {
		t.Fatal(err)
	}
	if err := event.AssociateOccurrence(incident.ID, "dokploy-app", "instance-a", "occurrence-1"); err != nil {
		t.Fatal(err)
	}
	if err := monitoring.CreateLogEvent(ctx, event); err != nil {
		t.Fatal(err)
	}

	incidents, _ := incidentrepo.New(db)
	investigations, _ := investigationrepo.New(db)
	operations, _ := operationrepo.New(db)
	evidences, _ := evidencerepo.New(db)
	projects, _ := projectrepo.New(db)
	issueReferences, _ := githubissuereferencerepo.New(db)
	transactor := dbadapter.NewTransactor(db)
	clock := monotonicClock(base)
	githubAccounts := fixedAccountReader{account: account}
	assembler := evidenceassembly.New(incidents, projects, githubAccounts, uuid.New, clock)
	runner := &evidenceAwareRunner{now: base.Add(10 * time.Second)}
	publisher := fixedIssuePublisher{number: 99, url: "https://github.com/acme/service-a/issues/99", createdAt: base.Add(20 * time.Second)}
	run := investigationusecase.NewRunUseCase(incidents, investigations, operations, evidences, projects, githubAccounts, issueReferences, publisher, issuecontent.New(), assembler, runner, transactor, clock)
	dispatcher := &recordingDispatcher{}
	start := investigationusecase.New(incidents, investigations, operations, transactor, dispatcher, uuid.New, clock)

	rollbackIncident, _ := domain.NewIncident(uuid.New(), "AKR-701", project.ID,
		"sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", domain.SeverityError, "rollback", base)
	if err := monitoring.CreateIncident(ctx, rollbackIncident); err != nil {
		t.Fatal(err)
	}
	rollbackStart := investigationusecase.New(incidents, investigations, failingCreateOperationStore{OperationStore: operations}, transactor, &recordingDispatcher{}, uuid.New, clock)
	if _, err := rollbackStart.StartIncidentInvestigation(ctx, portsin.StartIncidentInvestigationCommand{IncidentID: rollbackIncident.ID, IdempotencyKey: uuid.New()}); err == nil {
		t.Fatal("expected operation creation failure")
	}
	rolledBackIncident, err := incidents.Get(ctx, rollbackIncident.ID)
	if err != nil || rolledBackIncident.Phase != domain.IncidentPhaseDetected {
		t.Fatalf("Incident transition survived transaction rollback: %+v err=%v", rolledBackIncident, err)
	}
	rolledBackInvestigations, err := investigations.ListByIncident(ctx, rollbackIncident.ID, paging.Params{Limit: 25})
	if err != nil || rolledBackInvestigations.Total != 0 {
		t.Fatalf("Investigation survived transaction rollback: %+v err=%v", rolledBackInvestigations, err)
	}

	operation, err := start.StartIncidentInvestigation(ctx, portsin.StartIncidentInvestigationCommand{IncidentID: incident.ID, IdempotencyKey: uuid.New()})
	if err != nil {
		t.Fatal(err)
	}
	if operation.ResourceID == nil || dispatcher.investigationID != *operation.ResourceID {
		t.Fatalf("start did not commit/dispatch the durable pair: %+v %+v", operation, dispatcher)
	}
	if err := run.Execute(ctx, *operation.ResourceID, operation.ID); err != nil {
		t.Fatal(err)
	}

	completed, err := investigations.FindByID(ctx, *operation.ResourceID)
	if err != nil {
		t.Fatal(err)
	}
	if completed.Status != domain.InvestigationStatusCompleted || completed.RootCauseStatus == nil || *completed.RootCauseStatus != domain.RootCauseIdentified ||
		completed.ResolutionStatus == nil || *completed.ResolutionStatus != domain.ResolutionFixable || completed.EvidenceCount < 3 || len(completed.EvidenceIDs) != 2 {
		t.Fatalf("incomplete persisted structured result: %+v", completed)
	}
	storedIncident, err := incidents.Get(ctx, incident.ID)
	if err != nil || storedIncident.Phase != domain.IncidentPhasePublishingIssue {
		t.Fatalf("H4 did not leave fixable Incident waiting for remediation: incident=%+v err=%v", storedIncident, err)
	}
	reference, err := issueReferences.FindByInvestigation(ctx, completed.ID)
	if err != nil || reference == nil || reference.Number != 99 || reference.IncidentID != incident.ID {
		t.Fatalf("IssueReference was not persisted: reference=%+v err=%v", reference, err)
	}
	page, err := evidenceusecase.New(investigations, evidences).ListInvestigationEvidence(ctx, completed.ID, paging.Params{Limit: 25})
	if err != nil || len(page.Items) != completed.EvidenceCount {
		t.Fatalf("Evidence is not retrievable: page=%+v err=%v", page, err)
	}
	joined := ""
	for _, evidence := range page.Items {
		joined += evidence.Content + evidence.FilePath
	}
	for _, want := range []string{"database connection refused", "opening database connection", "retry scheduled", "internal/db/connect.go"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("persisted Evidence missing %q: %s", want, joined)
		}
	}
	if runner.context.Repository.Owner != "acme" || runner.context.Repository.Name != "service-a" {
		t.Fatalf("repository scope not resolved from Incident Project: %+v", runner.context.Repository)
	}

	if err := db.Exec("DELETE FROM investigations WHERE id = ?", completed.ID).Error; err == nil {
		t.Fatal("Evidence audit history did not RESTRICT Investigation deletion")
	}
	if err := db.Exec("DELETE FROM incidents WHERE id = ?", incident.ID).Error; err == nil {
		t.Fatal("Investigation audit history did not RESTRICT Incident deletion")
	}
	firstActive, _ := domain.NewInvestigation(uuid.New(), incident.ID, base.Add(time.Hour))
	secondActive, _ := domain.NewInvestigation(uuid.New(), incident.ID, base.Add(2*time.Hour))
	if err := investigations.Create(ctx, firstActive); err != nil {
		t.Fatal(err)
	}
	if err := investigations.Create(ctx, secondActive); !errors.Is(err, domain.ErrInvestigationAlreadyActive) {
		t.Fatalf("active unique index did not map conflict: %v", err)
	}
	if err := schema.SCHEMA_20260823_06_AddGitHubIssueReferences().Rollback(db); err != nil {
		t.Fatalf("rollback issue references: %v", err)
	}
	if err := schema.SCHEMA_20260823_05_AddInvestigationEvidenceIDs().Rollback(db); err != nil {
		t.Fatalf("rollback evidence_ids: %v", err)
	}
	if err := schema.SCHEMA_20260823_04_LinkInvestigationHistory().Rollback(db); err != nil {
		t.Fatalf("rollback history links: %v", err)
	}
	if db.Migrator().HasColumn("investigations", "evidence_ids") {
		t.Fatal("evidence_ids rollback did not remove the additive column")
	}
	if err := schema.SCHEMA_20260823_04_LinkInvestigationHistory().Migrate(db); err != nil {
		t.Fatalf("reapply history links: %v", err)
	}
	if err := schema.SCHEMA_20260823_05_AddInvestigationEvidenceIDs().Migrate(db); err != nil {
		t.Fatalf("reapply evidence_ids: %v", err)
	}
	if err := schema.SCHEMA_20260823_06_AddGitHubIssueReferences().Migrate(db); err != nil {
		t.Fatalf("reapply issue references: %v", err)
	}
}

func insertH3Project(t *testing.T, db *gorm.DB, now time.Time) (*domain.GitHubAccount, *domain.Project) {
	t.Helper()
	account, _ := domain.NewGitHubAccount(uuid.New(), "GitHub", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "acme", domain.IntegrationStatusConnected, now)
	server, _ := domain.NewDokployServer(uuid.New(), "Dokploy", "https://dokploy.example.com", "server", domain.IntegrationStatusConnected, now)
	if err := db.Table("github_accounts").Create(account).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Table("dokploy_servers").Create(server).Error; err != nil {
		t.Fatal(err)
	}
	repository, _ := domain.NewGitHubRepository(account.ID, "42", "acme", "service-a", "main", true, "https://github.com/acme/service-a")
	application, _ := domain.NewDokployApplication(server.ID, "dokploy-app", "instance-a", "Service A", "production", domain.DokployApplicationRunning)
	project, err := domain.NewProject(uuid.New(), "Service A", "", repository, application, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Table("projects").Create(project).Error; err != nil {
		t.Fatal(err)
	}
	return account, project
}

func monotonicClock(base time.Time) func() time.Time {
	var lock sync.Mutex
	step := 0
	return func() time.Time {
		lock.Lock()
		defer lock.Unlock()
		step++
		return base.Add(time.Duration(step) * time.Second)
	}
}
