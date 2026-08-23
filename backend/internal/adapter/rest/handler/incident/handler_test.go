package incident

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type incidentUseCaseStub struct {
	incident *domain.Incident
	events   []domain.LogEvent
}

func (s incidentUseCaseStub) Get(context.Context, uuid.UUID) (*domain.Incident, error) {
	if s.incident == nil {
		return nil, domain.ErrIncidentNotFound
	}
	return s.incident, nil
}

func (s incidentUseCaseStub) List(context.Context, paging.Params) (paging.Slice[domain.Incident], error) {
	if s.incident == nil {
		return paging.Slice[domain.Incident]{}, nil
	}
	return paging.Slice[domain.Incident]{Items: []domain.Incident{*s.incident}, Total: 1}, nil
}

func (s incidentUseCaseStub) ListLogEvents(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.LogEvent], error) {
	return paging.Slice[domain.LogEvent]{Items: s.events, Total: int64(len(s.events))}, nil
}

func TestIncidentListDetailAndLogEventsSerialization(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	projectID, incidentID := uuid.New(), uuid.New()
	value, err := domain.NewIncident(incidentID, "AKR-1", projectID, "sha256:"+strings.Repeat("a", 64), domain.SeverityError, "request failed", now)
	if err != nil {
		t.Fatal(err)
	}
	value.Project = &domain.ProjectReference{ID: projectID, Name: "payments"}
	record, _ := domain.NewSanitizedLogRecord(now.Add(-time.Second), domain.LogStreamUnknown, "before")
	event, err := domain.NewLogEvent(uuid.New(), projectID, now, domain.SeverityError, "ERROR request failed", value.Fingerprint, []string{"error_level"}, []domain.SanitizedLogRecord{record}, nil)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := New(incidentUseCaseStub{incident: value, events: []domain.LogEvent{*event}}, pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour})
	if err != nil {
		t.Fatal(err)
	}

	list := httptest.NewRecorder()
	handler.List(list, httptest.NewRequest(http.MethodGet, "/incidents", nil))
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), `"project":{"id":"`+projectID.String()+`","name":"payments"}`) || strings.Contains(list.Body.String(), "latest_investigation") {
		t.Fatalf("list response = %d %s", list.Code, list.Body.String())
	}

	detailRequest := httptest.NewRequest(http.MethodGet, "/incidents/"+incidentID.String(), nil)
	detailRequest.SetPathValue("incident_id", incidentID.String())
	detail := httptest.NewRecorder()
	handler.Get(detail, detailRequest)
	if detail.Code != http.StatusOK || !strings.Contains(detail.Body.String(), `"occurrence_count":1`) {
		t.Fatalf("detail response = %d %s", detail.Code, detail.Body.String())
	}

	eventsRequest := httptest.NewRequest(http.MethodGet, "/incidents/"+incidentID.String()+"/log-events", nil)
	eventsRequest.SetPathValue("incident_id", incidentID.String())
	events := httptest.NewRecorder()
	handler.ListLogEvents(events, eventsRequest)
	if events.Code != http.StatusOK || !strings.Contains(events.Body.String(), `"raw_context_redacted":true`) || !strings.Contains(events.Body.String(), `"message":"before"`) {
		t.Fatalf("events response = %d %s", events.Code, events.Body.String())
	}
}

func TestIncidentListRejectsInvalidPublishedFilters(t *testing.T) {
	handler, _ := New(incidentUseCaseStub{}, pagination.Config{Secret: []byte("01234567890123456789012345678901"), TTL: time.Hour})
	response := httptest.NewRecorder()
	handler.List(response, httptest.NewRequest(http.MethodGet, "/incidents?severity_in=bogus", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", response.Code)
	}
}
