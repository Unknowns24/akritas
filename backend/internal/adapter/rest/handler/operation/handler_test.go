package operation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type operationUseCaseStub struct {
	result *domain.Operation
	err    error
}

func (s *operationUseCaseStub) GetOperation(ctx context.Context, id uuid.UUID) (*domain.Operation, error) {
	return s.result, s.err
}

func TestGetHappyPath(t *testing.T) {
	useCase := &operationUseCaseStub{result: &domain.Operation{
		ID: uuid.New(), Type: domain.OperationTypeInvestigation, Status: domain.OperationStatusSucceeded,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}}
	handler, err := New(useCase)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/operations/"+useCase.result.ID.String(), nil)
	request.SetPathValue("operation_id", useCase.result.ID.String())
	recorder := httptest.NewRecorder()
	handler.Get(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestGetRejectsInvalidID(t *testing.T) {
	useCase := &operationUseCaseStub{}
	handler, err := New(useCase)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/operations/not-a-uuid", nil)
	request.SetPathValue("operation_id", "not-a-uuid")
	recorder := httptest.NewRecorder()
	handler.Get(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestGetPropagatesNotFound(t *testing.T) {
	useCase := &operationUseCaseStub{err: domain.ErrOperationNotFound}
	handler, err := New(useCase)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/operations/"+uuid.NewString(), nil)
	request.SetPathValue("operation_id", uuid.NewString())
	recorder := httptest.NewRecorder()
	handler.Get(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}
