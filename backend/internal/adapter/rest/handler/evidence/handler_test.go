package evidence

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type evidenceUseCaseStub struct {
	listResult          paging.Slice[domain.Evidence]
	listErr             error
	listInvestigationID uuid.UUID
	listCalled          bool
}

func (s *evidenceUseCaseStub) ListInvestigationEvidence(ctx context.Context, investigationID uuid.UUID, params paging.Params) (paging.Slice[domain.Evidence], error) {
	s.listCalled = true
	s.listInvestigationID = investigationID
	return s.listResult, s.listErr
}

func newHandlerFixture(t *testing.T) (*Handler, *evidenceUseCaseStub) {
	t.Helper()
	useCase := &evidenceUseCaseStub{}
	handler, err := New(useCase, pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	return handler, useCase
}

func TestListRejectsInvalidInvestigationID(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	request := httptest.NewRequest(http.MethodGet, "/investigations/not-a-uuid/evidence", nil)
	request.SetPathValue("investigation_id", "not-a-uuid")
	recorder := httptest.NewRecorder()
	handler.List(recorder, request)
	if recorder.Code != http.StatusBadRequest || useCase.listCalled {
		t.Fatalf("status/called = %d/%v, want 400/false", recorder.Code, useCase.listCalled)
	}
}

func TestListRejectsInvalidTypeIn(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	investigationID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, "/investigations/"+investigationID.String()+"/evidence?type_in=not_a_type", nil)
	request.SetPathValue("investigation_id", investigationID.String())
	recorder := httptest.NewRecorder()
	handler.List(recorder, request)
	if recorder.Code != http.StatusBadRequest || useCase.listCalled {
		t.Fatalf("status/called = %d/%v, want 400/false", recorder.Code, useCase.listCalled)
	}
}

func TestListAcceptsValidTypeIn(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	investigationID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, "/investigations/"+investigationID.String()+"/evidence?type_in=deployment_metadata,commit", nil)
	request.SetPathValue("investigation_id", investigationID.String())
	recorder := httptest.NewRecorder()
	handler.List(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
	if !useCase.listCalled || useCase.listInvestigationID != investigationID {
		t.Fatal("expected the list to be scoped to the given investigation")
	}
}

func TestListPropagatesNotFound(t *testing.T) {
	handler, useCase := newHandlerFixture(t)
	useCase.listErr = domain.ErrInvestigationNotFound
	investigationID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, "/investigations/"+investigationID.String()+"/evidence", nil)
	request.SetPathValue("investigation_id", investigationID.String())
	recorder := httptest.NewRecorder()
	handler.List(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}
