package monitoring

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/detection"
)

func sanitizeRecords(records []portsout.RawLogRecord) []domain.SanitizedLogRecord {
	result := make([]domain.SanitizedLogRecord, 0, len(records))
	for _, record := range records {
		message, _ := detection.Sanitize(record.Message)
		value, err := domain.NewSanitizedLogRecord(record.Timestamp, domain.LogStreamUnknown, message)
		if err == nil {
			result = append(result, value)
		}
	}
	return result
}

func buildDrafts(project domain.Project, checkpoint *domain.MonitoringCheckpoint, raw []portsout.RawLogRecord, now time.Time) ([]domain.PendingLogOccurrence, domain.MonitoringAssemblyState) {
	configuration := project.MonitoringConfiguration
	state := checkpoint.State
	records := sanitizeRecords(raw)
	ready := make([]domain.PendingLogOccurrence, 0)

	for index := range state.Pending {
		pending := &state.Pending[index]
		needed := pending.AfterRequired - len(pending.ContextAfter)
		if needed > len(records) {
			needed = len(records)
		}
		if needed > 0 {
			pending.ContextAfter = append(pending.ContextAfter, records[:needed]...)
		}
	}
	remainingPending := make([]domain.PendingLogOccurrence, 0, len(state.Pending))
	for _, pending := range state.Pending {
		if len(pending.ContextAfter) >= pending.AfterRequired || !now.Before(pending.FinalizeAfter) {
			ready = append(ready, pending)
		} else {
			remainingPending = append(remainingPending, pending)
		}
	}
	state.Pending = remainingPending

	engine, err := detection.NewEngine(configuration)
	if err == nil && configuration.Enabled {
		events := detection.Reconstruct(records)
		physicalOffset := 0
		for _, event := range events {
			start := physicalOffset
			physicalOffset += len(event.Records)
			detected := engine.Detect(project.ID, event)
			if detected == nil {
				continue
			}
			beforePool := append(append([]domain.SanitizedLogRecord(nil), state.RecentRecords...), records[:start]...)
			before := tail(beforePool, configuration.ContextBefore)
			after := head(records[physicalOffset:], configuration.ContextAfter)
			identity := detected.Timestamp.Format(time.RFC3339Nano) + "\x00" + detected.Message
			if start < len(raw) {
				identity = detected.Timestamp.Format(time.RFC3339Nano) + "\x00" + itoa(raw[start].Ordinal) + "\x00" + raw[start].ContentHash
			}
			pending := domain.PendingLogOccurrence{
				OccurrenceKey: occurrenceKey(checkpoint, identity), Timestamp: detected.Timestamp,
				Severity: detected.Severity, Message: detected.Message, Fingerprint: detected.Fingerprint,
				DetectionRules: append([]string(nil), detected.Rules...), ContextBefore: before, ContextAfter: after,
				AfterRequired: configuration.ContextAfter, FinalizeAfter: now.Add(pendingTimeout),
			}
			if len(after) >= configuration.ContextAfter {
				ready = append(ready, pending)
			} else {
				state.Pending = append(state.Pending, pending)
			}
		}
	}
	state.RecentRecords = tail(append(state.RecentRecords, records...), configuration.ContextBefore)
	state.OpenRecords = []domain.SanitizedLogRecord{}
	return ready, state
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var buffer [20]byte
	position := len(buffer)
	for value > 0 {
		position--
		buffer[position] = byte('0' + value%10)
		value /= 10
	}
	return string(buffer[position:])
}

func occurrenceKey(checkpoint *domain.MonitoringCheckpoint, identity string) string {
	payload := checkpoint.ID.String() + "\x00" + checkpoint.SourceApplicationID + "\x00" + checkpoint.SourceInstanceID + "\x00" + identity
	sum := sha256.Sum256([]byte(payload))
	return "occurrence:" + hex.EncodeToString(sum[:])
}

func title(message string) string {
	line := strings.TrimSpace(strings.SplitN(message, "\n", 2)[0])
	if len(line) > 500 {
		line = line[:500]
	}
	return line
}

func head(values []domain.SanitizedLogRecord, count int) []domain.SanitizedLogRecord {
	if count <= 0 {
		return []domain.SanitizedLogRecord{}
	}
	if len(values) > count {
		values = values[:count]
	}
	return append([]domain.SanitizedLogRecord(nil), values...)
}

func tail(values []domain.SanitizedLogRecord, count int) []domain.SanitizedLogRecord {
	if count <= 0 {
		return []domain.SanitizedLogRecord{}
	}
	if len(values) > count {
		values = values[len(values)-count:]
	}
	return append([]domain.SanitizedLogRecord(nil), values...)
}
