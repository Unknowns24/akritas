package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type InitialLogIngestion string

const (
	InitialLogIngestionFromNow   InitialLogIngestion = "from_now"
	InitialLogIngestionLast10000 InitialLogIngestion = "last_10000"
)

func (v InitialLogIngestion) Validate() error {
	switch v {
	case InitialLogIngestionFromNow, InitialLogIngestionLast10000:
		return nil
	default:
		return ErrInvalidInitialLogIngestion.Wrap(validationCause("initial log ingestion"))
	}
}

type MonitoringCursor struct {
	Timestamp   time.Time `json:"timestamp"`
	Ordinal     int       `json:"ordinal"`
	ContentHash string    `json:"content_hash"`
}

func (c MonitoringCursor) IsZero() bool { return c.Timestamp.IsZero() }

type PendingLogOccurrence struct {
	OccurrenceKey  string               `json:"occurrence_key"`
	Timestamp      time.Time            `json:"timestamp"`
	Severity       Severity             `json:"severity"`
	Message        string               `json:"message"`
	Fingerprint    string               `json:"fingerprint"`
	DetectionRules []string             `json:"detection_rules"`
	ContextBefore  []SanitizedLogRecord `json:"context_before"`
	ContextAfter   []SanitizedLogRecord `json:"context_after"`
	AfterRequired  int                  `json:"after_required"`
	FinalizeAfter  time.Time            `json:"finalize_after"`
}

type MonitoringAssemblyState struct {
	RecentRecords []SanitizedLogRecord   `json:"recent_records"`
	OpenRecords   []SanitizedLogRecord   `json:"open_records"`
	Pending       []PendingLogOccurrence `json:"pending"`
}

type MonitoringCheckpoint struct {
	ID                     uuid.UUID               `gorm:"column:id;type:uuid;primaryKey"`
	ProjectID              uuid.UUID               `gorm:"column:project_id;type:uuid"`
	SourceType             DokploySourceType       `gorm:"column:source_type"`
	SourceResourceID       string                  `gorm:"column:source_resource_id"`
	SourceServiceName      string                  `gorm:"column:source_service_name"`
	SourceInstanceID       string                  `gorm:"column:source_instance_id"`
	IsCurrent              bool                    `gorm:"column:is_current"`
	InitialBackfillPending bool                    `gorm:"column:initial_backfill_pending"`
	CursorTimestamp        *time.Time              `gorm:"column:cursor_timestamp"`
	CursorOrdinal          int                     `gorm:"column:cursor_ordinal"`
	CursorContentHash      string                  `gorm:"column:cursor_content_hash"`
	Version                int64                   `gorm:"column:version"`
	State                  MonitoringAssemblyState `gorm:"serializer:json;type:jsonb;column:assembly_state"`
	NextFinalizeAt         *time.Time              `gorm:"column:next_finalize_at"`
	CreatedAt              time.Time               `gorm:"column:created_at"`
	UpdatedAt              time.Time               `gorm:"column:updated_at"`
}

func NewMonitoringCheckpoint(id uuid.UUID, project Project, ingestion InitialLogIngestion, now time.Time) (*MonitoringCheckpoint, error) {
	if ingestion == "" {
		ingestion = InitialLogIngestionFromNow
	}
	if id == uuid.Nil || project.ID == uuid.Nil || ingestion.Validate() != nil || now.IsZero() {
		return nil, ErrInvalidMonitoringConfiguration.Wrap(validationCause("monitoring checkpoint"))
	}
	checkpoint := &MonitoringCheckpoint{
		ID: id, ProjectID: project.ID, SourceType: project.DokploySource.Type,
		SourceResourceID: project.DokploySource.ResourceIdentifier, SourceServiceName: project.DokploySource.ServiceName,
		SourceInstanceID: project.DokploySource.InstanceIdentifier, IsCurrent: true,
		InitialBackfillPending: ingestion == InitialLogIngestionLast10000, Version: 1,
		State:     MonitoringAssemblyState{RecentRecords: []SanitizedLogRecord{}, OpenRecords: []SanitizedLogRecord{}, Pending: []PendingLogOccurrence{}},
		CreatedAt: now.UTC(), UpdatedAt: now.UTC(),
	}
	if ingestion == InitialLogIngestionFromNow && project.MonitoringConfiguration.Enabled {
		anchor := now.UTC()
		checkpoint.CursorTimestamp = &anchor
		checkpoint.CursorContentHash = "anchor"
	}
	return checkpoint, nil
}

func (c MonitoringCheckpoint) Cursor() MonitoringCursor {
	if c.CursorTimestamp == nil {
		return MonitoringCursor{}
	}
	return MonitoringCursor{Timestamp: c.CursorTimestamp.UTC(), Ordinal: c.CursorOrdinal, ContentHash: c.CursorContentHash}
}

func (c *MonitoringCheckpoint) Advance(cursor MonitoringCursor, now time.Time) {
	timestamp := cursor.Timestamp.UTC()
	c.CursorTimestamp = &timestamp
	c.CursorOrdinal = cursor.Ordinal
	c.CursorContentHash = strings.TrimSpace(cursor.ContentHash)
	c.InitialBackfillPending = false
	c.Version++
	c.UpdatedAt = now.UTC()
}
