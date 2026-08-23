package evidenceassembly

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type fakeIncidentEvidenceReader struct {
	getResult *domain.Incident
	getErr    error
	logs      []domain.LogEvent
}

func (f *fakeIncidentEvidenceReader) Get(context.Context, uuid.UUID) (*domain.Incident, error) {
	return f.getResult, f.getErr
}
func (f *fakeIncidentEvidenceReader) ListLogEvents(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.LogEvent], error) {
	return paging.Slice[domain.LogEvent]{Items: append([]domain.LogEvent(nil), f.logs...), Total: int64(len(f.logs))}, nil
}

type fakeProjectStore struct {
	getResult *domain.Project
	getErr    error
}

func (f *fakeProjectStore) Create(context.Context, *domain.Project) error { return nil }
func (f *fakeProjectStore) Get(context.Context, uuid.UUID) (*domain.Project, error) {
	return f.getResult, f.getErr
}
func (f *fakeProjectStore) FindByNormalizedName(context.Context, string) (*domain.Project, error) {
	return nil, nil
}
func (f *fakeProjectStore) FindByDokploySource(context.Context, domain.DokploySourceSelector) (*domain.Project, error) {
	return nil, nil
}
func (f *fakeProjectStore) List(context.Context, paging.Params) (paging.Slice[domain.Project], error) {
	return paging.Slice[domain.Project]{}, nil
}
func (f *fakeProjectStore) Update(context.Context, *domain.Project, time.Time) error { return nil }
func (f *fakeProjectStore) Delete(context.Context, uuid.UUID) error                  { return nil }

type fakeAccountReader struct{ account *domain.GitHubAccount }

func (f fakeAccountReader) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return f.account, nil
}

type fakeCommitReader struct {
	commits []portsout.RepositoryCommitSummary
	err     error
}

func (f fakeCommitReader) ListRecentCommits(context.Context, domain.GitHubAccount, string, string, string, int) ([]portsout.RepositoryCommitSummary, error) {
	return append([]portsout.RepositoryCommitSummary(nil), f.commits...), f.err
}

func fixtureProject(t *testing.T) domain.Project {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	repository, err := domain.NewGitHubRepository(uuid.New(), "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	application, err := domain.NewDokployApplication(uuid.New(), "app-1", "api", "API", "production", domain.DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	source, err := domain.SourceFromApplication(application)
	if err != nil {
		t.Fatal(err)
	}
	project, err := domain.NewProject(uuid.New(), "Akritas", "demo", repository, source, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	return *project
}

func fixtureLogEvent(t *testing.T, incidentID, projectID uuid.UUID, now time.Time) domain.LogEvent {
	t.Helper()
	before, err := domain.NewSanitizedLogRecord(now.Add(-time.Second), domain.LogStreamStderr, "query started")
	if err != nil {
		t.Fatal(err)
	}
	after, err := domain.NewSanitizedLogRecord(now.Add(time.Second), domain.LogStreamStderr, "retry scheduled")
	if err != nil {
		t.Fatal(err)
	}
	event, err := domain.NewLogEvent(uuid.New(), projectID, now, domain.SeverityError,
		"database connection refused\nTraceback File app.py:42", "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		[]string{string(domain.DetectionRuleStackTrace)}, []domain.SanitizedLogRecord{before}, []domain.SanitizedLogRecord{after})
	if err != nil {
		t.Fatal(err)
	}
	if err := event.AssociateOccurrence(incidentID, evidenceSource(t, "app-1", "instance-1"), "occ-1"); err != nil {
		t.Fatal(err)
	}
	return *event
}

func TestAssembleProducesRealH2EvidenceAndFixedRepositoryScope(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	project := fixtureProject(t)
	incident := &domain.Incident{ID: uuid.New(), ProjectID: project.ID, Title: "DB error", Severity: domain.SeverityError}
	incidents := &fakeIncidentEvidenceReader{getResult: incident, logs: []domain.LogEvent{fixtureLogEvent(t, incident.ID, project.ID, now)}}
	assembler := New(incidents, &fakeProjectStore{getResult: &project}, fakeAccountReader{account: &domain.GitHubAccount{ID: project.GitHubRepository.GitHubAccountID}}, uuid.New, func() time.Time { return now })

	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), IncidentID: incident.ID, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	if result.Repository.Owner != "Unknowns24" || result.Repository.Name != "akritas" || result.Repository.Branch != "main" {
		t.Fatalf("unexpected fixed repository scope: %+v", result.Repository)
	}
	if len(result.Evidence) != 3 {
		t.Fatalf("expected deployment, log and real stack evidence, got %+v", result.Evidence)
	}
	joined := result.Evidence[1].Content
	for _, want := range []string{"database connection refused", "query started", "retry scheduled", "app-1", "instance-1", "stack_trace"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("log Evidence missing %q: %s", want, joined)
		}
	}
	if result.Evidence[0].Type != domain.EvidenceDeploymentMetadata || result.Evidence[1].Type != domain.EvidenceLogExcerpt || result.Evidence[2].Type != domain.EvidenceStackTrace {
		t.Fatalf("unexpected evidence types: %+v", result.Evidence)
	}
}

func TestAssembleDoesNotInventStackEvidence(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	project := fixtureProject(t)
	incident := &domain.Incident{ID: uuid.New(), ProjectID: project.ID}
	event := fixtureLogEvent(t, incident.ID, project.ID, now)
	event.DetectionRules = []string{string(domain.DetectionRuleErrorLevel)}
	assembler := New(&fakeIncidentEvidenceReader{getResult: incident, logs: []domain.LogEvent{event}}, &fakeProjectStore{getResult: &project}, fakeAccountReader{account: &domain.GitHubAccount{ID: project.GitHubRepository.GitHubAccountID}}, uuid.New, func() time.Time { return now })
	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), IncidentID: incident.ID, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	for _, evidence := range result.Evidence {
		if evidence.Type == domain.EvidenceStackTrace || evidence.Type == domain.EvidenceCommit || evidence.Type == domain.EvidenceDiff {
			t.Fatalf("invented Evidence: %+v", evidence)
		}
	}
}

func TestAssembleAddsSafePotentialCommitEvidence(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	project := fixtureProject(t)
	incident := &domain.Incident{
		ID: uuid.New(), ProjectID: project.ID, FirstSeenAt: now, LastSeenAt: now.Add(time.Minute),
		Title: "DB error", Severity: domain.SeverityError,
	}
	incidents := &fakeIncidentEvidenceReader{getResult: incident, logs: []domain.LogEvent{fixtureLogEvent(t, incident.ID, project.ID, now)}}
	commits := fakeCommitReader{commits: []portsout.RepositoryCommitSummary{{
		SHA: "deadbeef", Date: now.Add(-time.Hour).Format(time.RFC3339), Author: "dev",
		Message: "fix TOKEN=secret-value", URL: "https://github.com/acme/api/commit/deadbeef",
	}}}
	assembler := NewWithCommitCorrelation(incidents, &fakeProjectStore{getResult: &project}, fakeAccountReader{account: &domain.GitHubAccount{ID: project.GitHubRepository.GitHubAccountID}}, commits, uuid.New, func() time.Time { return now })

	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), IncidentID: incident.ID, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	var found *domain.Evidence
	for index := range result.Evidence {
		if result.Evidence[index].Type == domain.EvidenceCommit {
			found = &result.Evidence[index]
			break
		}
	}
	if found == nil {
		t.Fatalf("expected commit Evidence, got %+v", result.Evidence)
	}
	if found.CommitSHA != "deadbeef" || strings.Contains(found.Content, "secret-value") {
		t.Fatalf("unsafe commit Evidence: %+v", found)
	}
	if !strings.Contains(found.Summary, "potencialmente relacionado") || !strings.Contains(found.Summary, "no es causa confirmada") {
		t.Fatalf("commit Evidence must not claim causality: %s", found.Summary)
	}
}

func TestAssembleCommitCorrelationFailureDoesNotBlockInvestigation(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	project := fixtureProject(t)
	incident := &domain.Incident{ID: uuid.New(), ProjectID: project.ID, FirstSeenAt: now, LastSeenAt: now, Title: "DB error", Severity: domain.SeverityError}
	incidents := &fakeIncidentEvidenceReader{getResult: incident, logs: []domain.LogEvent{fixtureLogEvent(t, incident.ID, project.ID, now)}}
	assembler := NewWithCommitCorrelation(incidents, &fakeProjectStore{getResult: &project}, fakeAccountReader{account: &domain.GitHubAccount{ID: project.GitHubRepository.GitHubAccountID}}, fakeCommitReader{err: errors.New("github unavailable TOKEN=secret")}, uuid.New, func() time.Time { return now })

	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), IncidentID: incident.ID, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	for _, evidence := range result.Evidence {
		if evidence.Type == domain.EvidenceCommit || strings.Contains(evidence.Content, "TOKEN=secret") {
			t.Fatalf("commit correlation failure leaked or blocked incorrectly: %+v", result.Evidence)
		}
	}
}

func TestAssemblePropagatesMissingIncidentAndInfrastructureErrors(t *testing.T) {
	t.Parallel()
	want := errors.New("database unavailable")
	assembler := New(&fakeIncidentEvidenceReader{getErr: want}, &fakeProjectStore{}, fakeAccountReader{}, uuid.New, time.Now)
	_, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), IncidentID: uuid.New(), CreatedAt: time.Now()})
	if !errors.Is(err, want) {
		t.Fatalf("expected error to propagate, got %v", err)
	}
}

func TestAssembleEnforcesInitialCorpusCountAndRedactsSecrets(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	project := fixtureProject(t)
	incident := &domain.Incident{ID: uuid.New(), ProjectID: project.ID}
	logs := make([]domain.LogEvent, 0, 30)
	for index := 0; index < 30; index++ {
		event, err := domain.NewLogEvent(uuid.New(), project.ID, now.Add(time.Duration(index)*time.Second), domain.SeverityError,
			"connection failed token=very-secret-value", "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			[]string{string(domain.DetectionRuleErrorLevel)}, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := event.AssociateOccurrence(incident.ID, evidenceSource(t, "app", "instance"), uuid.NewString()); err != nil {
			t.Fatal(err)
		}
		logs = append(logs, *event)
	}
	assembler := New(&fakeIncidentEvidenceReader{getResult: incident, logs: logs}, &fakeProjectStore{getResult: &project}, fakeAccountReader{account: &domain.GitHubAccount{ID: project.GitHubRepository.GitHubAccountID}}, uuid.New, func() time.Time { return now })
	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), IncidentID: incident.ID, CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Evidence) > maximumInitialEvidence {
		t.Fatalf("initial Evidence count=%d", len(result.Evidence))
	}
	total := 0
	for _, evidence := range result.Evidence {
		total += evidenceSize(evidence)
		if strings.Contains(evidence.Content, "very-secret-value") {
			t.Fatal("secret leaked into persisted Evidence")
		}
	}
	if total > maximumCorpusBytes {
		t.Fatalf("initial persisted corpus=%d", total)
	}
}

func evidenceSource(t *testing.T, identifier, instance string) domain.DokploySource {
	t.Helper()
	source, err := domain.NewDokploySource(uuid.New(), domain.DokploySourceApplication, identifier, "", instance, "App", "production", domain.DokploySourceUnknown, "", "")
	if err != nil {
		t.Fatal(err)
	}
	return source
}
