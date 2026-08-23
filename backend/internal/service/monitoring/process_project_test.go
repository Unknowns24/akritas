package monitoring

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestDisabledProjectFinalizesStateWithoutAcquiringLogs(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	configuration := domain.DefaultMonitoringConfiguration()
	project := &domain.Project{ID: uuid.New(), MonitoringStatus: domain.MonitoringStatusDisabled, MonitoringConfiguration: configuration, DokploySource: domain.DokploySource{Type: domain.DokploySourceApplication, ResourceIdentifier: "app", InstanceIdentifier: "instance"}}
	checkpoint := &domain.MonitoringCheckpoint{ID: uuid.New(), ProjectID: project.ID, SourceType: domain.DokploySourceApplication, SourceResourceID: "app", SourceInstanceID: "instance", Version: 1, State: domain.MonitoringAssemblyState{Pending: []domain.PendingLogOccurrence{}}}
	store := &memoryMonitoringStore{project: project, checkpoint: checkpoint}
	servers := &serverReaderStub{}
	logs := &logSourceStub{}
	service, err := New(store, servers, logs, immediateTransactor{}, nil, nil, uuid.New, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	if err := service.ProcessProject(context.Background(), *project); err != nil {
		t.Fatal(err)
	}
	if logs.calls != 0 || servers.calls != 0 {
		t.Fatalf("disabled Project acquired provider state: logs=%d servers=%d", logs.calls, servers.calls)
	}
	if store.checkpoint.Version != 2 || project.MonitoringStatus != domain.MonitoringStatusDisabled {
		t.Fatalf("state was not finalized durably: checkpoint=%+v project=%+v", store.checkpoint, project)
	}
}

func TestPersistenceFailureDoesNotAdvanceCursor(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	configuration := domain.DefaultMonitoringConfiguration()
	configuration.Enabled = true
	serverID := uuid.New()
	project := &domain.Project{ID: uuid.New(), MonitoringStatus: domain.MonitoringStatusStarting, MonitoringConfiguration: configuration, DokploySource: domain.DokploySource{Type: domain.DokploySourceApplication, DokployServerID: serverID, ResourceIdentifier: "app", InstanceIdentifier: "instance"}}
	anchor := now.Add(-time.Minute)
	checkpoint := &domain.MonitoringCheckpoint{ID: uuid.New(), ProjectID: project.ID, SourceType: domain.DokploySourceApplication, SourceResourceID: "app", SourceInstanceID: "instance", Version: 1, CursorTimestamp: &anchor, CursorContentHash: "anchor", State: domain.MonitoringAssemblyState{}}
	persistenceFailure := errors.New("persistence failed")
	store := &memoryMonitoringStore{project: project, checkpoint: checkpoint, updateCheckpointErr: persistenceFailure}
	servers := &serverReaderStub{}
	logs := &logSourceStub{records: []portsout.RawLogRecord{{Timestamp: now, Ordinal: 0, ContentHash: "new", Message: "ordinary log"}}}
	service, err := New(store, servers, logs, immediateTransactor{}, nil, nil, uuid.New, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	if err := service.ProcessProject(context.Background(), *project); !errors.Is(err, persistenceFailure) {
		t.Fatalf("error = %v", err)
	}
	if store.checkpoint.Version != 1 || store.checkpoint.CursorContentHash != "anchor" || !store.checkpoint.CursorTimestamp.Equal(anchor) {
		t.Fatalf("cursor advanced after failure: %+v", store.checkpoint)
	}
	if project.MonitoringStatus != domain.MonitoringStatusError {
		t.Fatalf("initial failure status = %s", project.MonitoringStatus)
	}
}

func TestShouldStartAutomaticInvestigationForNewDetectedOrFailedIncidents(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	incident, err := domain.NewIncident(uuid.New(), "AKR-1", uuid.New(), "fingerprint", domain.SeverityError, "failure", now)
	if err != nil {
		t.Fatal(err)
	}
	if !shouldStartAutomaticInvestigation(incident, true) {
		t.Fatal("created incident should start automatic investigation")
	}
	if !shouldStartAutomaticInvestigation(incident, false) {
		t.Fatal("existing detected incident should start automatic investigation")
	}
	if err := incident.StartInvestigation(); err != nil {
		t.Fatal(err)
	}
	if shouldStartAutomaticInvestigation(incident, false) {
		t.Fatal("investigating incident should not start another automatic investigation")
	}
	if err := incident.FailInvestigation(); err != nil {
		t.Fatal(err)
	}
	if !shouldStartAutomaticInvestigation(incident, false) {
		t.Fatal("failed incident should be eligible for retry")
	}
	if shouldStartAutomaticInvestigation(nil, true) {
		t.Fatal("nil incident should not start automatic investigation")
	}
}
