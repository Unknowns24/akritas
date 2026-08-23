package investigation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type investigationUseCaseStub struct {
	startCalls   int
	startCommand portsin.StartIncidentInvestigationCommand
	startResult  *domain.Operation
	startErr     error
	getResult    *domain.Investigation
	getErr       error
	listResult   paging.Slice[domain.Investigation]
	listErr      error
}

func (s *investigationUseCaseStub) StartIncidentInvestigation(ctx context.Context, command portsin.StartIncidentInvestigationCommand) (*domain.Operation, error) {
	s.startCalls++
	s.startCommand = command
	return s.startResult, s.startErr
}

func (s *investigationUseCaseStub) GetInvestigation(ctx context.Context, id uuid.UUID) (*domain.Investigation, error) {
	return s.getResult, s.getErr
}

func (s *investigationUseCaseStub) ListIncidentInvestigations(ctx context.Context, incidentID uuid.UUID, params paging.Params) (paging.Slice[domain.Investigation], error) {
	return s.listResult, s.listErr
}

func newHandlerFixture(t *testing.T) (*Handler, *investigationUseCaseStub) {
	t.Helper()
	useCase := &investigationUseCaseStub{}
	handler, err := New(useCase, pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	return handler, useCase
}

func TestStartRejectsMissingIdempotencyKey(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	request := httptest.NewRequest(http.MethodPost, "/incidents/"+uuid.NewString()+"/investigations", nil)
	request.SetPathValue("incident_id", uuid.NewString())
	recorder := httptest.NewRecorder()
	handler.Start(recorder, request)
	if recorder.Code != http.StatusBadRequest || useCase.startCalls != 0 {
		t.Fatalf("status/calls = %d/%d, want 400/0", recorder.Code, useCase.startCalls)
	}
}

func TestStartHappyPathReturns202WithOperation(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	incidentID := uuid.New()
	key := uuid.New()
	useCase.startResult = &domain.Operation{ID: uuid.New(), Type: domain.OperationTypeInvestigation, Status: domain.OperationStatusQueued, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	request := httptest.NewRequest(http.MethodPost, "/incidents/"+incidentID.String()+"/investigations", nil)
	request.SetPathValue("incident_id", incidentID.String())
	request.Header.Set("Idempotency-Key", key.String())
	recorder := httptest.NewRecorder()
	handler.Start(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
	if useCase.startCalls != 1 || useCase.startCommand.IncidentID != incidentID || useCase.startCommand.IdempotencyKey != key {
		t.Fatalf("unexpected command passed to the use case: %+v", useCase.startCommand)
	}
}

func TestStartPropagatesIncidentNotFound(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	useCase.startErr = domain.ErrIncidentNotFound
	request := httptest.NewRequest(http.MethodPost, "/incidents/"+uuid.NewString()+"/investigations", nil)
	request.SetPathValue("incident_id", uuid.NewString())
	request.Header.Set("Idempotency-Key", uuid.NewString())
	recorder := httptest.NewRecorder()
	handler.Start(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}

func TestGetRejectsInvalidID(t *testing.T) {
	handler, _ := newHandlerFixture(t)
	request := httptest.NewRequest(http.MethodGet, "/investigations/not-a-uuid", nil)
	request.SetPathValue("investigation_id", "not-a-uuid")
	recorder := httptest.NewRecorder()
	handler.Get(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestGetPropagatesNotFound(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	useCase.getErr = domain.ErrInvestigationNotFound
	request := httptest.NewRequest(http.MethodGet, "/investigations/"+uuid.NewString(), nil)
	request.SetPathValue("investigation_id", uuid.NewString())
	recorder := httptest.NewRecorder()
	handler.Get(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}

func TestListRejectsInvalidIncidentID(t *testing.T) {
	handler, _ := newHandlerFixture(t)
	request := httptest.NewRequest(http.MethodGet, "/incidents/not-a-uuid/investigations", nil)
	request.SetPathValue("incident_id", "not-a-uuid")
	recorder := httptest.NewRecorder()
	handler.List(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}
