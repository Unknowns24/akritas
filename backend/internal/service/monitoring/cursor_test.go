package monitoring

import (
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func TestRecordsAfterSkipsOverlappingCursorDeterministically(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	checkpoint := &domain.MonitoringCheckpoint{CursorTimestamp: &now, CursorOrdinal: 1, CursorContentHash: "b"}
	records := []portsout.RawLogRecord{
		{Timestamp: now, Ordinal: 0, ContentHash: "a"},
		{Timestamp: now, Ordinal: 1, ContentHash: "b"},
		{Timestamp: now.Add(time.Second), Ordinal: 0, ContentHash: "c"},
	}
	got, err := recordsAfter(checkpoint, records)
	if err != nil || len(got) != 1 || got[0].ContentHash != "c" {
		t.Fatalf("recordsAfter = %+v, %v", got, err)
	}
}

func TestRecordsAfterFailsClosedWhenSaturatedWindowCannotProveContinuity(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	checkpoint := &domain.MonitoringCheckpoint{CursorTimestamp: &now, CursorContentHash: "missing"}
	records := make([]portsout.RawLogRecord, maximumFetchRecords)
	for index := range records {
		records[index] = portsout.RawLogRecord{Timestamp: now.Add(-time.Second), Ordinal: index, ContentHash: "other"}
	}
	if _, err := recordsAfter(checkpoint, records); !errors.Is(err, domain.ErrMonitoringContinuityLost) {
		t.Fatalf("error = %v", err)
	}
}

func TestFromNowAnchorFailsClosedWhenMoreThanTailArrivedBeforeFirstPoll(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	checkpoint := &domain.MonitoringCheckpoint{CursorTimestamp: &now, CursorContentHash: "anchor"}
	records := make([]portsout.RawLogRecord, maximumFetchRecords)
	for index := range records {
		records[index] = portsout.RawLogRecord{Timestamp: now.Add(time.Second), Ordinal: index, ContentHash: "new"}
	}
	if _, err := recordsAfter(checkpoint, records); !errors.Is(err, domain.ErrMonitoringContinuityLost) {
		t.Fatalf("error = %v", err)
	}
}
