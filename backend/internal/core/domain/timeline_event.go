package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type TimelineEventType string

const (
	TimelineIncidentDetected     TimelineEventType = "incident_detected"
	TimelineInvestigationStarted TimelineEventType = "investigation_started"
	TimelineToolUsed             TimelineEventType = "tool_used"
	TimelineRootCauseClassified  TimelineEventType = "root_cause_classified"
	TimelineIssueCreated         TimelineEventType = "issue_created"
	TimelineWorkflowFailed       TimelineEventType = "workflow_failed"
)

func (t TimelineEventType) Validate() error {
	switch t {
	case TimelineIncidentDetected, TimelineInvestigationStarted, TimelineToolUsed,
		TimelineRootCauseClassified, TimelineIssueCreated, TimelineWorkflowFailed:
		return nil
	default:
		return ErrInvalidIncident.Wrap(validationCause("timeline event type"))
	}
}

type TimelineEvent struct {
	ID         uuid.UUID         `gorm:"column:id"`
	IncidentID uuid.UUID         `gorm:"column:incident_id"`
	Type       TimelineEventType `gorm:"column:type"`
	OccurredAt time.Time         `gorm:"column:occurred_at"`
	Summary    string            `gorm:"column:summary"`
	Detail     string            `gorm:"column:detail"`
}

func (e TimelineEvent) Validate() error {
	if e.ID == uuid.Nil || e.IncidentID == uuid.Nil || e.Type.Validate() != nil ||
		!validTime(e.OccurredAt) || !nonBlank(e.Summary) || len(e.Summary) > 2000 || len(e.Detail) > 10000 {
		return ErrInvalidIncident.Wrap(validationCause("timeline event"))
	}
	e.Summary = strings.TrimSpace(e.Summary)
	e.Detail = strings.TrimSpace(e.Detail)
	return nil
}
