package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMonitoringCheckpointInitialIngestionSemantics(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	project := Project{ID: uuid.New(), DokploySource: DokploySource{Type: DokploySourceApplication, ResourceIdentifier: "app", InstanceIdentifier: "instance"}, MonitoringConfiguration: DefaultMonitoringConfiguration()}
	fromNow, err := NewMonitoringCheckpoint(uuid.New(), project, InitialLogIngestionFromNow, now)
	if err != nil || fromNow.CursorTimestamp != nil || fromNow.InitialBackfillPending {
		t.Fatalf("disabled from-now checkpoint = %+v, %v", fromNow, err)
	}
	backfill, err := NewMonitoringCheckpoint(uuid.New(), project, InitialLogIngestionLast10000, now)
	if err != nil || !backfill.InitialBackfillPending || backfill.CursorTimestamp != nil {
		t.Fatalf("backfill checkpoint = %+v, %v", backfill, err)
	}
	project.MonitoringConfiguration.Enabled = true
	enabled, err := NewMonitoringCheckpoint(uuid.New(), project, InitialLogIngestionFromNow, now)
	if err != nil || enabled.CursorTimestamp == nil || !enabled.CursorTimestamp.Equal(now) || enabled.CursorContentHash != "anchor" {
		t.Fatalf("enabled from-now checkpoint = %+v, %v", enabled, err)
	}
}
