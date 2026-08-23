package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type InvestigationStatus string

const (
	InvestigationStatusPending   InvestigationStatus = "pending"
	InvestigationStatusRunning   InvestigationStatus = "running"
	InvestigationStatusCompleted InvestigationStatus = "completed"
	InvestigationStatusFailed    InvestigationStatus = "failed"
)

func (s InvestigationStatus) Validate() error {
	switch s {
	case InvestigationStatusPending, InvestigationStatusRunning, InvestigationStatusCompleted, InvestigationStatusFailed:
		return nil
	default:
		return ErrInvalidInvestigationStatus.Wrap(validationCause("investigation status"))
	}
}

type Investigation struct {
	ID                 uuid.UUID           `gorm:"column:id;type:uuid;primaryKey"`
	IncidentID         uuid.UUID           `gorm:"column:incident_id;type:uuid"`
	Status             InvestigationStatus `gorm:"column:status"`
	CreatedAt          time.Time           `gorm:"column:created_at"`
	StartedAt          *time.Time          `gorm:"column:started_at"`
	FinishedAt         *time.Time          `gorm:"column:finished_at"`
	Summary            string              `gorm:"column:summary"`
	RootCause          string              `gorm:"column:root_cause"`
	RootCauseStatus    *RootCauseStatus    `gorm:"column:root_cause_status"`
	ResolutionStatus   *ResolutionStatus   `gorm:"column:resolution_status"`
	Confidence         *float64            `gorm:"column:confidence"`
	Hypotheses         []string            `gorm:"serializer:json;type:jsonb;column:hypotheses"`
	RelevantFiles      []string            `gorm:"serializer:json;type:jsonb;column:relevant_files"`
	RelevantCommits    []string            `gorm:"serializer:json;type:jsonb;column:relevant_commits"`
	RecommendedActions []string            `gorm:"serializer:json;type:jsonb;column:recommended_actions"`
	EvidenceCount      int                 `gorm:"column:evidence_count"`
	FailureUserMessage string              `gorm:"column:failure_user_message"`
}

func NewInvestigation(id, incidentID uuid.UUID, createdAt time.Time) (*Investigation, error) {
	investigation := &Investigation{
		ID: id, IncidentID: incidentID, Status: InvestigationStatusPending, CreatedAt: createdAt,
		Hypotheses: []string{}, RelevantFiles: []string{}, RelevantCommits: []string{}, RecommendedActions: []string{},
	}
	if err := investigation.Validate(); err != nil {
		return nil, err
	}
	return investigation, nil
}

func (i Investigation) Validate() error {
	if i.ID == uuid.Nil || i.IncidentID == uuid.Nil || i.Status.Validate() != nil || !validTime(i.CreatedAt) || i.EvidenceCount < 0 {
		return ErrInvalidInvestigation.Wrap(validationCause("investigation"))
	}
	switch i.Status {
	case InvestigationStatusPending:
		if i.StartedAt != nil || i.FinishedAt != nil {
			return ErrInvalidInvestigation.Wrap(validationCause("pending investigation times"))
		}
	case InvestigationStatusRunning:
		if i.StartedAt == nil || i.StartedAt.Before(i.CreatedAt) || i.FinishedAt != nil {
			return ErrInvalidInvestigation.Wrap(validationCause("running investigation times"))
		}
	case InvestigationStatusCompleted:
		if !i.validTerminalTimes() || !nonBlank(i.Summary) || i.RootCauseStatus == nil || i.ResolutionStatus == nil || i.Confidence == nil ||
			i.RootCauseStatus.Validate() != nil || i.ResolutionStatus.Validate() != nil || !validConfidence(*i.Confidence) {
			return ErrInvalidInvestigation.Wrap(validationCause("completed investigation"))
		}
	case InvestigationStatusFailed:
		if !i.validTerminalTimes() || !nonBlank(i.FailureUserMessage) {
			return ErrInvalidInvestigation.Wrap(validationCause("failed investigation"))
		}
	}
	return nil
}

func (i Investigation) validTerminalTimes() bool {
	return i.StartedAt != nil && i.FinishedAt != nil && !i.StartedAt.Before(i.CreatedAt) && !i.FinishedAt.Before(*i.StartedAt)
}

func (i *Investigation) Start(at time.Time) error {
	if i == nil || i.Status != InvestigationStatusPending || at.Before(i.CreatedAt) {
		return ErrInvestigationTransition.Wrap(validationCause("start investigation"))
	}
	i.Status = InvestigationStatusRunning
	i.StartedAt = &at
	return nil
}

func (i *Investigation) Complete(
	at time.Time,
	summary, rootCause string,
	rootCauseStatus RootCauseStatus,
	resolutionStatus ResolutionStatus,
	confidence float64,
	hypotheses, relevantFiles, relevantCommits, recommendedActions []string,
) error {
	if i == nil || i.Status != InvestigationStatusRunning || i.StartedAt == nil || at.Before(*i.StartedAt) ||
		!nonBlank(summary) || rootCauseStatus.Validate() != nil || resolutionStatus.Validate() != nil || !validConfidence(confidence) {
		return ErrInvestigationTransition.Wrap(validationCause("complete investigation"))
	}
	i.Status = InvestigationStatusCompleted
	i.FinishedAt = &at
	i.Summary = strings.TrimSpace(summary)
	i.RootCause = strings.TrimSpace(rootCause)
	i.RootCauseStatus = &rootCauseStatus
	i.ResolutionStatus = &resolutionStatus
	i.Confidence = &confidence
	i.Hypotheses = cloneStrings(hypotheses)
	i.RelevantFiles = cloneStrings(relevantFiles)
	i.RelevantCommits = cloneStrings(relevantCommits)
	i.RecommendedActions = cloneStrings(recommendedActions)
	return nil
}

func (i *Investigation) Fail(at time.Time, userMessage string) error {
	if i == nil || i.Status != InvestigationStatusRunning || i.StartedAt == nil || at.Before(*i.StartedAt) || !nonBlank(userMessage) {
		return ErrInvestigationTransition.Wrap(validationCause("fail investigation"))
	}
	i.Status = InvestigationStatusFailed
	i.FinishedAt = &at
	i.FailureUserMessage = strings.TrimSpace(userMessage)
	return nil
}
