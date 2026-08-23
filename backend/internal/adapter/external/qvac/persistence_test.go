package qvac_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/external/qvac"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	investigationusecase "github.com/Unknowns24/akritas/backend/internal/usecase/investigation"
	"github.com/google/uuid"
)

type memInvestigationStore struct {
	byID map[uuid.UUID]*domain.Investigation
}

func (s *memInvestigationStore) Create(ctx context.Context, value *domain.Investigation) error {
	s.byID[value.ID] = cloneInvestigation(value)
	return nil
}
func (s *memInvestigationStore) Update(ctx context.Context, value *domain.Investigation) error {
	s.byID[value.ID] = cloneInvestigation(value)
	return nil
}
func (s *memInvestigationStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Investigation, error) {
	value, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrInvestigationNotFound
	}
	return cloneInvestigation(value), nil
}
func (s *memInvestigationStore) FindLatestByIncident(ctx context.Context, incidentID uuid.UUID) (*domain.Investigation, error) {
	for _, value := range s.byID {
		if value.IncidentID == incidentID {
			return cloneInvestigation(value), nil
		}
	}
	return nil, domain.ErrInvestigationNotFound
}
func (s *memInvestigationStore) ListByIncident(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.Investigation], error) {
	return paging.Slice[domain.Investigation]{}, nil
}
func (s *memInvestigationStore) ExistsActiveForIncident(ctx context.Context, incidentID uuid.UUID) (bool, error) {
	return false, nil
}
func (s *memInvestigationStore) ListOpen(context.Context) ([]domain.Investigation, error) {
	return nil, nil
}

func cloneInvestigation(value *domain.Investigation) *domain.Investigation {
	cloned := *value
	return &cloned
}

type memOperationStore struct {
	byID map[uuid.UUID]*domain.Operation
}

func (s *memOperationStore) Create(ctx context.Context, value *domain.Operation) error {
	copied := *value
	s.byID[value.ID] = &copied
	return nil
}
func (s *memOperationStore) Update(ctx context.Context, value *domain.Operation) error {
	copied := *value
	s.byID[value.ID] = &copied
	return nil
}
func (s *memOperationStore) FindByID(ctx context.Context, id uuid.UUID) (*domain.Operation, error) {
	value, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrOperationNotFound
	}
	copied := *value
	return &copied, nil
}
func (s *memOperationStore) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Operation, error) {
	return nil, domain.ErrOperationNotFound
}
func (s *memOperationStore) FindByResource(context.Context, domain.OperationResourceType, uuid.UUID) (*domain.Operation, error) {
	return nil, domain.ErrOperationNotFound
}

type memEvidenceStore struct{ created []domain.Evidence }

func (s *memEvidenceStore) Create(ctx context.Context, value *domain.Evidence) error {
	s.created = append(s.created, *value)
	return nil
}
func (s *memEvidenceStore) ListByInvestigation(ctx context.Context, investigationID uuid.UUID, params paging.Params) (paging.Slice[domain.Evidence], error) {
	return paging.Slice[domain.Evidence]{Items: s.created}, nil
}

type emptyAssembler struct{}

func (emptyAssembler) Assemble(ctx context.Context, investigation domain.Investigation) (portsout.InvestigationRunContext, error) {
	return portsout.InvestigationRunContext{Investigation: investigation}, nil
}

type memIncidentStore struct{ incident *domain.Incident }

func (s *memIncidentStore) Get(context.Context, uuid.UUID) (*domain.Incident, error) {
	return s.incident, nil
}
func (s *memIncidentStore) Lock(context.Context, uuid.UUID) (*domain.Incident, error) {
	return s.incident, nil
}
func (s *memIncidentStore) Update(_ context.Context, value *domain.Incident) error {
	copied := *value
	s.incident = &copied
	return nil
}

type memProjectStore struct{ project *domain.Project }

func (s memProjectStore) Create(context.Context, *domain.Project) error { return nil }
func (s memProjectStore) Get(context.Context, uuid.UUID) (*domain.Project, error) {
	return s.project, nil
}
func (s memProjectStore) FindByNormalizedName(context.Context, string) (*domain.Project, error) {
	return nil, domain.ErrProjectNotFound
}
func (s memProjectStore) FindByDokployApplication(context.Context, uuid.UUID, string) (*domain.Project, error) {
	return nil, domain.ErrProjectNotFound
}
func (s memProjectStore) List(context.Context, paging.Params) (paging.Slice[domain.Project], error) {
	return paging.Slice[domain.Project]{}, nil
}
func (s memProjectStore) Update(context.Context, *domain.Project, time.Time) error { return nil }
func (s memProjectStore) Delete(context.Context, uuid.UUID) error                  { return nil }

type memGitHubAccountReader struct{ account *domain.GitHubAccount }

func (s memGitHubAccountReader) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return s.account, nil
}

type memIssueReferenceStore struct{ reference *domain.GitHubIssueReference }

func (s *memIssueReferenceStore) Create(ctx context.Context, value *domain.GitHubIssueReference) error {
	copied := *value
	s.reference = &copied
	return nil
}
func (s *memIssueReferenceStore) FindByInvestigation(context.Context, uuid.UUID) (*domain.GitHubIssueReference, error) {
	return nil, nil
}
func (s *memIssueReferenceStore) FindLatestByIncident(context.Context, uuid.UUID) (*domain.GitHubIssueReference, error) {
	return s.reference, nil
}

type memIssuePublisher struct{}

func (memIssuePublisher) PublishIssue(context.Context, domain.GitHubAccount, domain.GitHubRepository, portsout.IssueContent) (portsout.PublishedIssue, error) {
	return portsout.PublishedIssue{Number: 9, URL: "https://github.com/acme/service/issues/9", CreatedAt: time.Date(2026, 8, 23, 12, 10, 0, 0, time.UTC)}, nil
}

type memIssueContentBuilder struct{}

func (memIssueContentBuilder) BuildIssueContent(portsout.IssueContentInput) (portsout.IssueContent, error) {
	return portsout.IssueContent{Title: "issue", Body: "body"}, nil
}

type passthroughTransactor struct{}

func (passthroughTransactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestRunUseCasePersistsStructuredQVACResult(t *testing.T) {
	t.Parallel()
	payload := `{
		"summary":"worker panic",
		"root_cause":"nil map write",
		"root_cause_status":"identified",
		"resolution_status":"fixable",
		"confidence":0.88,
		"hypotheses":["race"],
		"evidence_ids":[],
		"relevant_files":["worker.go"],
		"relevant_commits":["abc123"],
		"recommended_actions":["initialize map"]
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": payload}}},
		})
	}))
	t.Cleanup(server.Close)

	client, err := qvac.NewClient(qvac.ClientConfig{EndpointURL: server.URL + "/v1", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := qvac.NewRunner(client, nil, qvac.RunnerConfig{})
	if err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	investigation, err := domain.NewInvestigation(uuid.New(), uuid.New(), now)
	if err != nil {
		t.Fatal(err)
	}
	operation, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, nil, nil, nil, "queued", now)
	if err != nil {
		t.Fatal(err)
	}
	investigations := &memInvestigationStore{byID: map[uuid.UUID]*domain.Investigation{investigation.ID: investigation}}
	operations := &memOperationStore{byID: map[uuid.UUID]*domain.Operation{operation.ID: operation}}
	account, _ := domain.NewGitHubAccount(uuid.New(), "Acme", domain.GitHubAccountPersonal, domain.GitHubAuthenticationPersonalAccessToken, "acme", domain.IntegrationStatusConnected, now)
	repository, _ := domain.NewGitHubRepository(account.ID, "42", "acme", "service", "main", true, "https://github.com/acme/service")
	application, _ := domain.NewDokployApplication(uuid.New(), "app", "instance", "Service", "prod", domain.DokployApplicationRunning)
	project, _ := domain.NewProject(uuid.New(), "Service", "", repository, application, domain.DefaultMonitoringConfiguration(), now)
	incident := &domain.Incident{ID: investigation.IncidentID, ProjectID: project.ID, Phase: domain.IncidentPhaseInvestigating}
	uc := investigationusecase.NewRunUseCase(
		&memIncidentStore{incident: incident}, investigations, operations, &memEvidenceStore{},
		memProjectStore{project: project}, memGitHubAccountReader{account: account}, &memIssueReferenceStore{},
		memIssuePublisher{}, memIssueContentBuilder{}, emptyAssembler{}, runner, passthroughTransactor{},
		func() time.Time { return now.Add(time.Second) },
	)
	if err := uc.Execute(context.Background(), investigation.ID, operation.ID); err != nil {
		t.Fatal(err)
	}
	stored := investigations.byID[investigation.ID]
	if stored.Status != domain.InvestigationStatusCompleted || stored.Summary != "worker panic" ||
		stored.RootCause != "nil map write" || stored.RootCauseStatus == nil || *stored.RootCauseStatus != domain.RootCauseIdentified ||
		stored.ResolutionStatus == nil || *stored.ResolutionStatus != domain.ResolutionFixable ||
		stored.Confidence == nil || *stored.Confidence != 0.88 ||
		len(stored.RelevantFiles) != 1 || stored.RelevantFiles[0] != "worker.go" ||
		len(stored.RecommendedActions) != 1 {
		t.Fatalf("persisted investigation = %+v", stored)
	}
	_ = portsout.InvestigationRunResult{}
}
