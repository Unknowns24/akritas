package dashboard

import (
	"time"

	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/google/uuid"
)

type ActivityEventDTO struct {
	ID          uuid.UUID                        `json:"id"`
	Type        string                           `json:"type"`
	Project     *incidentdto.ProjectReferenceDTO `json:"project,omitempty"`
	IncidentID  *uuid.UUID                       `json:"incident_id,omitempty"`
	IncidentKey string                           `json:"incident_key,omitempty"`
	OccurredAt  time.Time                        `json:"occurred_at"`
	Summary     string                           `json:"summary"`
	UserMessage string                           `json:"user_message,omitempty"`
}

type OverviewDTO struct {
	MonitoredProjects          int                       `json:"monitored_projects"`
	ActiveIncidents            int                       `json:"active_incidents"`
	WorkflowCompletedIncidents int                       `json:"workflow_completed_incidents"`
	PullRequestsCreated        int                       `json:"pull_requests_created"`
	ActiveInvestigations       []incidentdto.IncidentDTO `json:"active_investigations"`
	LatestActivity             []ActivityEventDTO        `json:"latest_activity,omitempty"`
}
