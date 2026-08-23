package evidenceassembly

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type fakeIncidentReader struct {
	getResult *domain.Incident
	getErr    error
}

func (f *fakeIncidentReader) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}
func (f *fakeIncidentReader) Get(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	return f.getResult, f.getErr
}

type fakeProjectStore struct {
	getResult *domain.Project
	getErr    error
}

func (f *fakeProjectStore) Create(ctx context.Context, value *domain.Project) error { return nil }
func (f *fakeProjectStore) Get(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return f.getResult, f.getErr
}
func (f *fakeProjectStore) FindByNormalizedName(ctx context.Context, name string) (*domain.Project, error) {
	return nil, nil
}
func (f *fakeProjectStore) FindByDokployApplication(ctx context.Context, serverID uuid.UUID, applicationIdentifier string) (*domain.Project, error) {
	return nil, nil
}
func (f *fakeProjectStore) List(ctx context.Context, params paging.Params) (paging.Slice[domain.Project], error) {
	return paging.Slice[domain.Project]{}, nil
}
func (f *fakeProjectStore) Update(ctx context.Context, value *domain.Project, expected time.Time) error {
	return nil
}
func (f *fakeProjectStore) Delete(ctx context.Context, id uuid.UUID) error { return nil }

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
	project, err := domain.NewProject(uuid.New(), "Akritas", "demo", repository, application, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	return *project
}

func TestAssembleProducesDeploymentMetadataFromRealProject(t *testing.T) {
	t.Parallel()
	project := fixtureProject(t)
	incident := &domain.Incident{ID: uuid.New(), ProjectID: uuid.New()}
	incidents := &fakeIncidentReader{getResult: incident}
	projects := &fakeProjectStore{getResult: &project}
	assembler := New(incidents, projects, uuid.New, func() time.Time { return time.Now() })

	investigation := domain.Investigation{ID: uuid.New(), IncidentID: incident.ID, CreatedAt: time.Now()}
	result, err := assembler.Assemble(context.Background(), investigation)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].Type != domain.EvidenceDeploymentMetadata {
		t.Fatalf("expected exactly one deployment_metadata Evidence, got %+v", result)
	}
	if !strings.Contains(result[0].Content, project.Name) || !strings.Contains(result[0].Content, project.GitHubRepository.FullName) {
		t.Fatalf("expected the content to reflect real project fields, got %s", result[0].Content)
	}
	if !result[0].Redacted {
		t.Fatal("assembled evidence must always be redacted")
	}
}

func TestAssembleReturnsEmptyWhenIncidentNotFound(t *testing.T) {
	t.Parallel()
	incidents := &fakeIncidentReader{getErr: domain.ErrIncidentNotFound}
	projects := &fakeProjectStore{}
	assembler := New(incidents, projects, uuid.New, func() time.Time { return time.Now() })

	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), CreatedAt: time.Now()})
	if err != nil || len(result) != 0 {
		t.Fatalf("expected no evidence and no error, got %+v, %v", result, err)
	}
}

func TestAssembleReturnsEmptyWhenProjectNotFound(t *testing.T) {
	t.Parallel()
	incidents := &fakeIncidentReader{getResult: &domain.Incident{ID: uuid.New(), ProjectID: uuid.New()}}
	projects := &fakeProjectStore{getErr: domain.ErrProjectNotFound}
	assembler := New(incidents, projects, uuid.New, func() time.Time { return time.Now() })

	result, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), CreatedAt: time.Now()})
	if err != nil || len(result) != 0 {
		t.Fatalf("expected no evidence and no error, got %+v, %v", result, err)
	}
}

func TestAssemblePropagatesInfrastructureErrors(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("database unavailable")
	incidents := &fakeIncidentReader{getErr: wantErr}
	projects := &fakeProjectStore{}
	assembler := New(incidents, projects, uuid.New, func() time.Time { return time.Now() })

	_, err := assembler.Assemble(context.Background(), domain.Investigation{ID: uuid.New(), CreatedAt: time.Now()})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected the infrastructure error to propagate, got %v", err)
	}
}
