package monitoring

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestPersistOccurrenceCreatesGroupsByLastSeenAndStartsNewWindow(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	store := &memoryMonitoringStore{}
	service := &Service{store: store, newID: uuid.New}
	project := domain.Project{
		ID: uuid.New(), DokploySource: domain.DokploySource{Type: domain.DokploySourceApplication, DokployServerID: uuid.New(), ResourceIdentifier: "app", InstanceIdentifier: "instance", DisplayName: "App", Status: domain.DokploySourceUnknown},
		MonitoringConfiguration: domain.MonitoringConfiguration{GroupingWindow: 30 * time.Minute},
	}
	occurrences := []domain.PendingLogOccurrence{
		occurrence("one", now, domain.SeverityError),
		occurrence("two", now.Add(20*time.Minute), domain.SeverityCritical),
		occurrence("three", now.Add(51*time.Minute), domain.SeverityError),
	}
	for index := range occurrences {
		if _, _, err := service.persistOccurrence(context.Background(), project, &occurrences[index]); err != nil {
			t.Fatal(err)
		}
	}
	if len(store.incidents) != 2 || len(store.events) != 3 {
		t.Fatalf("incidents=%d events=%d", len(store.incidents), len(store.events))
	}
	first := store.incidents[0]
	if first.OccurrenceCount != 2 || !first.FirstSeenAt.Equal(now) || !first.LastSeenAt.Equal(now.Add(20*time.Minute)) || first.Severity != domain.SeverityCritical {
		t.Fatalf("grouped incident = %+v", first)
	}
	if store.incidents[1].OccurrenceCount != 1 || !store.incidents[1].FirstSeenAt.Equal(now.Add(51*time.Minute)) {
		t.Fatalf("new incident = %+v", store.incidents[1])
	}

	duplicate := occurrence("two", now.Add(20*time.Minute), domain.SeverityCritical)
	if _, _, err := service.persistOccurrence(context.Background(), project, &duplicate); err != nil {
		t.Fatal(err)
	}
	if len(store.events) != 3 || store.incidents[0].OccurrenceCount != 2 {
		t.Fatal("retry duplicated durable effects")
	}
}

func TestAutomaticInvestigationQueuesCreatedIncidentWhenEnabled(t *testing.T) {
	project := domain.Project{ID: uuid.New(), Name: "payments"}
	incidentID := uuid.New()
	automation := &automationPolicyStub{policy: domain.AutomationPolicy{AutomaticInvestigation: true}}
	starter := &investigationStarterStub{}
	service := &Service{automation: automation, investigations: starter}

	service.startAutomaticInvestigations(context.Background(), project, []automaticInvestigationCandidate{{incidentID: incidentID, incidentKey: "AKR-1"}})

	if automation.calls != 1 {
		t.Fatalf("policy calls = %d", automation.calls)
	}
	if len(starter.commands) != 1 || starter.commands[0].IncidentID != incidentID || starter.commands[0].IdempotencyKey == uuid.Nil {
		t.Fatalf("start commands = %+v", starter.commands)
	}
	service.startAutomaticInvestigations(context.Background(), project, []automaticInvestigationCandidate{{incidentID: incidentID, incidentKey: "AKR-1"}})
	if len(starter.commands) != 2 || starter.commands[0].IdempotencyKey != starter.commands[1].IdempotencyKey {
		t.Fatalf("automatic idempotency key must be deterministic: %+v", starter.commands)
	}
}

func TestAutomaticInvestigationSkipsWhenDisabled(t *testing.T) {
	project := domain.Project{ID: uuid.New(), Name: "payments"}
	automation := &automationPolicyStub{policy: domain.AutomationPolicy{AutomaticInvestigation: false}}
	starter := &investigationStarterStub{}
	service := &Service{automation: automation, investigations: starter}

	service.startAutomaticInvestigations(context.Background(), project, []automaticInvestigationCandidate{{incidentID: uuid.New(), incidentKey: "AKR-1"}})

	if automation.calls != 1 {
		t.Fatalf("policy calls = %d", automation.calls)
	}
	if len(starter.commands) != 0 {
		t.Fatalf("unexpected start commands = %+v", starter.commands)
	}
}

func TestAutomaticInvestigationFailureDoesNotPanicOrRetryInline(t *testing.T) {
	project := domain.Project{ID: uuid.New(), Name: "payments"}
	starterFailure := errors.New("starter failed")
	automation := &automationPolicyStub{policy: domain.AutomationPolicy{AutomaticInvestigation: true}}
	starter := &investigationStarterStub{err: starterFailure}
	service := &Service{automation: automation, investigations: starter}

	service.startAutomaticInvestigations(context.Background(), project, []automaticInvestigationCandidate{{incidentID: uuid.New(), incidentKey: "AKR-1"}})

	if len(starter.commands) != 1 {
		t.Fatalf("start commands = %+v", starter.commands)
	}
}

func occurrence(key string, timestamp time.Time, severity domain.Severity) domain.PendingLogOccurrence {
	return domain.PendingLogOccurrence{OccurrenceKey: key, Timestamp: timestamp, Severity: severity, Message: "ERROR failure", Fingerprint: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", DetectionRules: []string{"error_level"}}
}

type memoryMonitoringStore struct {
	incidents           []*domain.Incident
	events              []*domain.LogEvent
	project             *domain.Project
	checkpoint          *domain.MonitoringCheckpoint
	updateCheckpointErr error
}

func (s *memoryMonitoringStore) NextIncidentKey(context.Context) (string, error) {
	return fmt.Sprintf("AKR-%d", len(s.incidents)+1), nil
}

func (s *memoryMonitoringStore) FindGroupableIncident(_ context.Context, projectID uuid.UUID, fingerprint string, occurredAt time.Time, window time.Duration) (*domain.Incident, error) {
	for index := len(s.incidents) - 1; index >= 0; index-- {
		if s.incidents[index].CanGroup(projectID, fingerprint, occurredAt, window) {
			return s.incidents[index], nil
		}
	}
	return nil, nil
}

func (s *memoryMonitoringStore) CreateIncident(_ context.Context, incident *domain.Incident) error {
	s.incidents = append(s.incidents, incident)
	return nil
}

func (*memoryMonitoringStore) UpdateIncident(context.Context, *domain.Incident) error { return nil }

func (s *memoryMonitoringStore) OccurrenceExists(_ context.Context, _ uuid.UUID, key string) (bool, error) {
	for _, event := range s.events {
		if event.OccurrenceKey == key {
			return true, nil
		}
	}
	return false, nil
}

func (s *memoryMonitoringStore) CreateLogEvent(_ context.Context, event *domain.LogEvent) error {
	s.events = append(s.events, event)
	return nil
}

func (*memoryMonitoringStore) ListProjectsForMonitoring(context.Context) ([]domain.Project, error) {
	return nil, nil
}
func (s *memoryMonitoringStore) LockProject(context.Context, uuid.UUID) (*domain.Project, error) {
	return s.project, nil
}
func (s *memoryMonitoringStore) GetCurrentCheckpoint(context.Context, uuid.UUID, bool) (*domain.MonitoringCheckpoint, error) {
	if s.checkpoint == nil {
		return nil, nil
	}
	copy := *s.checkpoint
	return &copy, nil
}

func (s *memoryMonitoringStore) CreateCheckpoint(_ context.Context, checkpoint *domain.MonitoringCheckpoint) error {
	s.checkpoint = checkpoint
	return nil
}
func (s *memoryMonitoringStore) RotateCheckpoint(_ context.Context, checkpoint *domain.MonitoringCheckpoint) error {
	s.checkpoint = checkpoint
	return nil
}
func (s *memoryMonitoringStore) UpdateCheckpoint(_ context.Context, checkpoint *domain.MonitoringCheckpoint, _ int64) error {
	if s.updateCheckpointErr != nil {
		return s.updateCheckpointErr
	}
	copy := *checkpoint
	s.checkpoint = &copy
	return nil
}
func (s *memoryMonitoringStore) UpdateProjectObservation(_ context.Context, _ uuid.UUID, status domain.MonitoringStatus, _ time.Time, _ *time.Time) error {
	if s.project != nil {
		s.project.MonitoringStatus = status
	}
	return nil
}

type immediateTransactor struct{}

func (immediateTransactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type serverReaderStub struct{ calls int }

func (s *serverReaderStub) Get(context.Context, uuid.UUID) (*domain.DokployServer, error) {
	s.calls++
	return &domain.DokployServer{}, nil
}

type logSourceStub struct {
	calls   int
	records []portsout.RawLogRecord
}

func (s *logSourceStub) FetchLogs(context.Context, portsout.LogFetchRequest) ([]portsout.RawLogRecord, error) {
	s.calls++
	return append([]portsout.RawLogRecord(nil), s.records...), nil
}

type automationPolicyStub struct {
	calls  int
	policy domain.AutomationPolicy
	err    error
}

func (s *automationPolicyStub) Get(context.Context) (domain.AutomationPolicy, error) {
	s.calls++
	if s.err != nil {
		return domain.AutomationPolicy{}, s.err
	}
	return s.policy, nil
}

type investigationStarterStub struct {
	commands []portsin.StartIncidentInvestigationCommand
	err      error
}

func (s *investigationStarterStub) StartIncidentInvestigation(_ context.Context, command portsin.StartIncidentInvestigationCommand) (*domain.Operation, error) {
	s.commands = append(s.commands, command)
	if s.err != nil {
		return nil, s.err
	}
	now := time.Now().UTC()
	operation, err := domain.NewOperation(uuid.New(), domain.OperationTypeInvestigation, nil, nil, nil, "queued", now)
	if err != nil {
		return nil, err
	}
	return operation, nil
}
