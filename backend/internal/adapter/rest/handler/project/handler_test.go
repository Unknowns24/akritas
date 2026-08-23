package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func TestCreateReturnsSafeStableProjectEnvelope(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	project := projectFixture(t, 0)
	useCase.createResult = &portsin.ProjectResult{Project: &project, BuiltInDetectionRules: domain.AllBuiltInDetectionRules()}
	body := fmt.Sprintf(`{
		"name":"Akritas","description":"demo","github_account_id":%q,
		"repository_identifier":"42","default_branch":"main","dokploy_server_id":%q,
		"application_identifier":"app-1","monitoring_configuration":{
			"enabled":false,"error_patterns":["panic"],"ignored_patterns":[],
			"grouping_window":"PT30M","context_before":20,"context_after":20}}
	`, project.GitHubRepository.GitHubAccountID, project.DokployApplication.DokployServerID)

	recorder := httptest.NewRecorder()
	handler.Create(recorder, httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(body)))
	if recorder.Code != http.StatusCreated || useCase.createCalls != 1 {
		t.Fatalf("Create status/calls = %d/%d, body=%s", recorder.Code, useCase.createCalls, recorder.Body.String())
	}
	responseBody := strings.ToLower(recorder.Body.String())
	for _, forbidden := range []string{"password", "secret", "api_key", "credential", "access_token"} {
		if strings.Contains(responseBody, forbidden) {
			t.Fatalf("Project response leaks forbidden name %q: %s", forbidden, recorder.Body.String())
		}
	}
	var envelope struct {
		Data struct {
			ID                      string `json:"id"`
			Name                    string `json:"name"`
			GitHubRepository        any    `json:"github_repository"`
			DokployApplication      any    `json:"dokploy_application"`
			MonitoringConfiguration any    `json:"monitoring_configuration"`
			BuiltInDetectionRules   []any  `json:"built_in_detection_rules"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Data.ID != project.ID.String() || envelope.Data.Name != "Akritas" || envelope.Data.GitHubRepository == nil || envelope.Data.DokployApplication == nil || envelope.Data.MonitoringConfiguration == nil || len(envelope.Data.BuiltInDetectionRules) != 7 {
		t.Fatalf("incomplete Project response: %+v", envelope.Data)
	}
}

func TestProjectHandlersRejectInvalidInputAndReturnStableConflicts(t *testing.T) {
	handler, useCase := newHandlerFixture(t)

	t.Run("invalid create body", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler.Create(recorder, httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(`{"name":"Akritas"}`)))
		if recorder.Code != http.StatusBadRequest || useCase.createCalls != 0 {
			t.Fatalf("status/calls = %d/%d", recorder.Code, useCase.createCalls)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/projects/not-a-uuid", nil)
		request.SetPathValue("project_id", "not-a-uuid")
		handler.Delete(recorder, request)
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", recorder.Code)
		}
	})

	t.Run("active delete conflict", func(t *testing.T) {
		useCase.deleteErr = domain.ErrProjectMustBeDisabled
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/projects/"+uuid.NewString(), nil)
		request.SetPathValue("project_id", uuid.NewString())
		handler.Delete(recorder, request)
		if recorder.Code != http.StatusConflict || !strings.Contains(recorder.Body.String(), domain.ErrProjectMustBeDisabled.Code) {
			t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("invalid monitoring regex", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/projects/id/monitoring-configuration", strings.NewReader(`{"enabled":true,"error_patterns":["["],"ignored_patterns":[],"grouping_window":"PT30M","context_before":20,"context_after":20}`))
		request.SetPathValue("project_id", uuid.NewString())
		handler.PutMonitoring(recorder, request)
		if recorder.Code != http.StatusBadRequest || useCase.putMonitoringCalls != 0 {
			t.Fatalf("status/calls = %d/%d", recorder.Code, useCase.putMonitoringCalls)
		}
	})

	t.Run("invalid list filter", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler.List(recorder, httptest.NewRequest(http.MethodGet, "/projects?monitoring_status_in=bogus", nil))
		if recorder.Code != http.StatusBadRequest || useCase.listCalls != 0 {
			t.Fatalf("status/calls = %d/%d", recorder.Code, useCase.listCalls)
		}
	})
}

func TestDeleteReturnsNoContentAndListBuildsSignedCursor(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	id := uuid.New()
	deleteRecorder := httptest.NewRecorder()
	deleteRequest := httptest.NewRequest(http.MethodDelete, "/projects/"+id.String(), nil)
	deleteRequest.SetPathValue("project_id", id.String())
	handler.Delete(deleteRecorder, deleteRequest)
	if deleteRecorder.Code != http.StatusNoContent || deleteRecorder.Body.Len() != 0 || useCase.deleteID != id {
		t.Fatalf("delete = %d %q id=%s", deleteRecorder.Code, deleteRecorder.Body.String(), useCase.deleteID)
	}

	useCase.listResult.Items = make([]domain.Project, 26)
	for index := range useCase.listResult.Items {
		useCase.listResult.Items[index] = projectFixture(t, index)
	}
	useCase.listResult.Total = 26
	listRecorder := httptest.NewRecorder()
	handler.List(listRecorder, httptest.NewRequest(http.MethodGet, "/projects?limit=25&name_like=Akritas", nil))
	if listRecorder.Code != http.StatusOK || useCase.listParams.Filters["name_like"] != "Akritas" {
		t.Fatalf("list = %d %s params=%+v", listRecorder.Code, listRecorder.Body.String(), useCase.listParams)
	}
	var envelope struct {
		Data   []any `json:"data"`
		Paging struct {
			NextCursor string `json:"next_cursor"`
		} `json:"paging"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if len(envelope.Data) != 25 || envelope.Paging.NextCursor == "" {
		t.Fatalf("page data/cursor = %d/%q", len(envelope.Data), envelope.Paging.NextCursor)
	}
}

func newHandlerFixture(t *testing.T) (*Handler, *projectUseCaseStub) {
	t.Helper()
	useCase := &projectUseCaseStub{}
	handler, err := New(useCase, pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	return handler, useCase
}

func projectFixture(t *testing.T, offset int) domain.Project {
	t.Helper()
	now := time.Date(2026, 8, 22, 12, offset, 0, 0, time.UTC)
	repository, err := domain.NewGitHubRepository(uuid.New(), fmt.Sprintf("%d", offset+1), "Unknowns24", fmt.Sprintf("akritas-%d", offset), "main", true, fmt.Sprintf("https://github.com/Unknowns24/akritas-%d", offset))
	if err != nil {
		t.Fatal(err)
	}
	application, err := domain.NewDokployApplication(uuid.New(), fmt.Sprintf("app-%d", offset), fmt.Sprintf("instance-%d", offset), "API", "production", domain.DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	value, err := domain.NewProject(uuid.New(), "Akritas", "demo", repository, application, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	return *value
}

type projectUseCaseStub struct {
	createResult       *portsin.ProjectResult
	createErr          error
	createCalls        int
	deleteErr          error
	deleteID           uuid.UUID
	putMonitoringCalls int
	listResult         paging.Slice[domain.Project]
	listParams         paging.Params
	listCalls          int
}

func (s *projectUseCaseStub) Create(context.Context, portsin.CreateProjectCommand) (*portsin.ProjectResult, error) {
	s.createCalls++
	return s.createResult, s.createErr
}
func (*projectUseCaseStub) Get(context.Context, uuid.UUID) (*portsin.ProjectResult, error) {
	return nil, domain.ErrProjectNotFound
}
func (s *projectUseCaseStub) List(_ context.Context, params paging.Params) (paging.Slice[domain.Project], error) {
	s.listCalls++
	s.listParams = params
	return s.listResult, nil
}
func (*projectUseCaseStub) Update(context.Context, portsin.UpdateProjectCommand) (*portsin.ProjectResult, error) {
	return nil, errors.New("unexpected update")
}
func (s *projectUseCaseStub) Delete(_ context.Context, id uuid.UUID) error {
	s.deleteID = id
	return s.deleteErr
}
func (*projectUseCaseStub) GetMonitoring(context.Context, uuid.UUID) (domain.MonitoringConfiguration, error) {
	return domain.MonitoringConfiguration{}, errors.New("unexpected get monitoring")
}
func (s *projectUseCaseStub) PutMonitoring(context.Context, uuid.UUID, domain.MonitoringConfiguration) (domain.MonitoringConfiguration, error) {
	s.putMonitoringCalls++
	return domain.MonitoringConfiguration{}, nil
}
