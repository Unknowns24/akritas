package dokploy

import (
	"context"
	"encoding/json"
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

func TestListComposeServicesValidatesRefreshAndMapsEnvelope(t *testing.T) {
	serverID := uuid.New()
	stub := &dokployUseCaseStub{}
	stub.services = []domain.DokployComposeService{
		{DokployServerID: serverID, ComposeIdentifier: "compose-1", ServiceName: "api"},
		{DokployServerID: serverID, ComposeIdentifier: "compose-1", ServiceName: "worker"},
	}
	handler, err := New(stub, pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/services?refresh=true", nil)
	request.SetPathValue("server_id", serverID.String())
	request.SetPathValue("compose_id", "compose-1")
	recorder := httptest.NewRecorder()
	handler.ListComposeServices(recorder, request)
	if recorder.Code != http.StatusOK || !stub.refresh || stub.composeID != "compose-1" {
		t.Fatalf("response=%d %s stub=%+v", recorder.Code, recorder.Body.String(), stub)
	}
	var envelope struct {
		Data []struct {
			ServiceName string `json:"service_name"`
		} `json:"data"`
	}
	if json.Unmarshal(recorder.Body.Bytes(), &envelope) != nil || len(envelope.Data) != 2 || envelope.Data[0].ServiceName != "api" {
		t.Fatalf("envelope = %+v", envelope)
	}

	invalid := httptest.NewRequest(http.MethodGet, "/services?refresh=1", nil)
	invalid.SetPathValue("server_id", serverID.String())
	invalid.SetPathValue("compose_id", "compose-1")
	invalidRecorder := httptest.NewRecorder()
	handler.ListComposeServices(invalidRecorder, invalid)
	if invalidRecorder.Code != http.StatusBadRequest || stub.calls != 1 {
		t.Fatalf("invalid refresh response=%d calls=%d", invalidRecorder.Code, stub.calls)
	}

	empty := httptest.NewRequest(http.MethodGet, "/services?refresh=", nil)
	empty.SetPathValue("server_id", serverID.String())
	empty.SetPathValue("compose_id", "compose-1")
	emptyRecorder := httptest.NewRecorder()
	handler.ListComposeServices(emptyRecorder, empty)
	if emptyRecorder.Code != http.StatusBadRequest || stub.calls != 1 {
		t.Fatalf("empty refresh response=%d calls=%d", emptyRecorder.Code, stub.calls)
	}
}

type dokployUseCaseStub struct {
	services  []domain.DokployComposeService
	refresh   bool
	composeID string
	calls     int
}

func (*dokployUseCaseStub) Create(context.Context, portsin.CreateDokployServerCommand) (*domain.DokployServer, error) {
	return nil, nil
}
func (*dokployUseCaseStub) Get(context.Context, uuid.UUID) (*domain.DokployServer, error) {
	return nil, nil
}
func (*dokployUseCaseStub) List(context.Context, paging.Params) (paging.Slice[domain.DokployServer], error) {
	return paging.Slice[domain.DokployServer]{}, nil
}
func (*dokployUseCaseStub) Update(context.Context, uuid.UUID, portsin.UpdateDokployServerCommand) (*domain.DokployServer, error) {
	return nil, nil
}
func (*dokployUseCaseStub) Delete(context.Context, uuid.UUID) error { return nil }
func (*dokployUseCaseStub) TestConnection(context.Context, uuid.UUID) (portsin.ConnectionTestResult, error) {
	return portsin.ConnectionTestResult{}, nil
}
func (*dokployUseCaseStub) ListApplications(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.DokployApplication], error) {
	return paging.Slice[domain.DokployApplication]{}, nil
}
func (*dokployUseCaseStub) ListComposes(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.DokployCompose], error) {
	return paging.Slice[domain.DokployCompose]{}, nil
}
func (s *dokployUseCaseStub) ListComposeServices(_ context.Context, _ uuid.UUID, composeID string, refresh bool) ([]domain.DokployComposeService, error) {
	s.calls++
	s.composeID = composeID
	s.refresh = refresh
	return s.services, nil
}
