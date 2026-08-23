package domain

import (
	"time"

	"github.com/google/uuid"
)

type ComponentHealthStatus string

const (
	ComponentHealthHealthy       ComponentHealthStatus = "healthy"
	ComponentHealthDegraded      ComponentHealthStatus = "degraded"
	ComponentHealthUnavailable   ComponentHealthStatus = "unavailable"
	ComponentHealthNotConfigured ComponentHealthStatus = "not_configured"
	ComponentHealthRunning       ComponentHealthStatus = "running"
)

type ComponentHealth struct {
	Component string
	Status    ComponentHealthStatus
	CheckedAt *time.Time
	Message   string
}

type SystemStatus struct {
	GitHubAccountCount int
	DokployServerCount int
	QvacEndpoint       string
	Components         []ComponentHealth
	LastDiagnosticsAt  *time.Time
}

type Overview struct {
	MonitoredProjects          int
	ActiveIncidents            int
	WorkflowCompletedIncidents int
	PullRequestsCreated        int
	ActiveInvestigations       []Incident
	LatestActivity             []ActivityEvent
}

type ActivityType string

const (
	ActivityTypeMonitoring    ActivityType = "monitoring"
	ActivityTypeIncident      ActivityType = "incident"
	ActivityTypeInvestigation ActivityType = "investigation"
	ActivityTypeIssue         ActivityType = "issue"
	ActivityTypeRemediation   ActivityType = "remediation"
	ActivityTypePullRequest   ActivityType = "pull_request"
	ActivityTypeIntegration   ActivityType = "integration"
	ActivityTypeSystem        ActivityType = "system"
)

func (t ActivityType) Validate() error {
	switch t {
	case ActivityTypeMonitoring, ActivityTypeIncident, ActivityTypeInvestigation, ActivityTypeIssue, ActivityTypeRemediation, ActivityTypePullRequest, ActivityTypeIntegration, ActivityTypeSystem:
		return nil
	default:
		return ErrInvalidOperationType.Wrap(validationCause("activity type"))
	}
}

type ActivityEvent struct {
	ID          uuid.UUID
	Type        ActivityType
	Project     *ProjectReference
	IncidentID  *uuid.UUID
	IncidentKey string
	OccurredAt  time.Time
	Summary     string
	UserMessage string
}

type PullRequestProjection struct {
	ID                uuid.UUID
	Project           ProjectReference
	IncidentID        uuid.UUID
	IncidentKey       string
	RemediationID     uuid.UUID
	IssueReference    *GitHubIssueReference
	Reference         PullRequestReference
	Title             string
	ChangesSummary    string
	ValidationSummary string
	RiskSummary       string
	CreatedAt         time.Time
}
