package monitoring

import (
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestBuildDraftsPersistsBoundedBeforeAndDelayedAfterContext(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	project := monitoringProject(t, now)
	checkpoint := &domain.MonitoringCheckpoint{
		ID: uuid.New(), ProjectID: project.ID, SourceApplicationID: "app", SourceInstanceID: "instance",
		State: domain.MonitoringAssemblyState{RecentRecords: []domain.SanitizedLogRecord{sanitized(t, now.Add(-2*time.Second), "before-1"), sanitized(t, now.Add(-time.Second), "before-2")}},
	}
	first := []portsout.RawLogRecord{
		raw(now, 0, "h0", "ordinary-before"), raw(now.Add(time.Second), 0, "h1", "BOOM request failed"), raw(now.Add(2*time.Second), 0, "h2", "after-1"),
	}
	ready, state := buildDrafts(project, checkpoint, first, now)
	if len(ready) != 0 || len(state.Pending) != 1 {
		t.Fatalf("ready=%d pending=%d", len(ready), len(state.Pending))
	}
	pending := state.Pending[0]
	if len(pending.ContextBefore) != 2 || pending.ContextBefore[0].Message != "before-2" || pending.ContextBefore[1].Message != "ordinary-before" || len(pending.ContextAfter) != 1 {
		t.Fatalf("pending context = %+v / %+v", pending.ContextBefore, pending.ContextAfter)
	}
	checkpoint.State = state
	ready, state = buildDrafts(project, checkpoint, []portsout.RawLogRecord{raw(now.Add(3*time.Second), 0, "h3", "after-2")}, now.Add(10*time.Second))
	if len(ready) != 1 || len(state.Pending) != 0 || len(ready[0].ContextAfter) != 2 {
		t.Fatalf("ready=%+v pending=%+v", ready, state.Pending)
	}
}

func TestBuildDraftsFinalizesPendingAtThirtySecondDeadline(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	project := monitoringProject(t, now)
	checkpoint := &domain.MonitoringCheckpoint{State: domain.MonitoringAssemblyState{Pending: []domain.PendingLogOccurrence{{OccurrenceKey: "key", AfterRequired: 2, FinalizeAfter: now.Add(pendingTimeout)}}}}
	ready, state := buildDrafts(project, checkpoint, nil, now.Add(pendingTimeout))
	if len(ready) != 1 || len(state.Pending) != 0 {
		t.Fatalf("ready=%+v pending=%+v", ready, state.Pending)
	}
}

func monitoringProject(t *testing.T, now time.Time) domain.Project {
	t.Helper()
	configuration, err := domain.NewMonitoringConfiguration(true, []string{"BOOM"}, []string{}, time.Minute, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	return domain.Project{ID: uuid.New(), MonitoringConfiguration: configuration}
}

func sanitized(t *testing.T, timestamp time.Time, message string) domain.SanitizedLogRecord {
	t.Helper()
	value, err := domain.NewSanitizedLogRecord(timestamp, domain.LogStreamUnknown, message)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func raw(timestamp time.Time, ordinal int, hash, message string) portsout.RawLogRecord {
	return portsout.RawLogRecord{Timestamp: timestamp, Ordinal: ordinal, ContentHash: hash, Message: message}
}
