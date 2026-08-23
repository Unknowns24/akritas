package monitoring

import (
	"strconv"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func fetchSince(checkpoint *domain.MonitoringCheckpoint, now time.Time) string {
	if checkpoint.InitialBackfillPending || checkpoint.CursorTimestamp == nil {
		return "all"
	}
	duration := now.Sub(checkpoint.CursorTimestamp.UTC()) + overlapDuration
	seconds := int64(duration / time.Second)
	if duration%time.Second != 0 {
		seconds++
	}
	if seconds < 1 {
		seconds = 1
	}
	return strconv.FormatInt(seconds, 10) + "s"
}

func recordsAfter(checkpoint *domain.MonitoringCheckpoint, records []portsout.RawLogRecord) ([]portsout.RawLogRecord, error) {
	if checkpoint.InitialBackfillPending || checkpoint.CursorTimestamp == nil {
		return records, nil
	}
	cursor := checkpoint.Cursor()
	if cursor.ContentHash == "anchor" {
		if len(records) == maximumFetchRecords && len(records) > 0 && records[0].Timestamp.After(cursor.Timestamp) {
			return nil, domain.ErrMonitoringContinuityLost
		}
		index := 0
		for index < len(records) && !records[index].Timestamp.After(cursor.Timestamp) {
			index++
		}
		return records[index:], nil
	}
	for index, record := range records {
		if record.Timestamp.Equal(cursor.Timestamp) && record.Ordinal == cursor.Ordinal && record.ContentHash == cursor.ContentHash {
			return records[index+1:], nil
		}
	}
	if len(records) == 0 {
		return nil, nil
	}
	if len(records) < maximumFetchRecords {
		return recordsStrictlyAfter(cursor.Timestamp, records), nil
	}
	return nil, domain.ErrMonitoringContinuityLost
}

func recordsStrictlyAfter(timestamp time.Time, records []portsout.RawLogRecord) []portsout.RawLogRecord {
	index := 0
	for index < len(records) && !records[index].Timestamp.After(timestamp) {
		index++
	}
	return records[index:]
}
