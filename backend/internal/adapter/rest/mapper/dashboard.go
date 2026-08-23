package mapper

import (
	dashboarddto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dashboard"
	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func OverviewToDTO(value domain.Overview) dashboarddto.OverviewDTO {
	incidents := make([]incidentdto.IncidentDTO, 0, len(value.ActiveInvestigations))
	for _, item := range value.ActiveInvestigations {
		incidents = append(incidents, IncidentToDTO(item))
	}
	activity := make([]dashboarddto.ActivityEventDTO, 0, len(value.LatestActivity))
	for _, item := range value.LatestActivity {
		activity = append(activity, ActivityEventToDTO(item))
	}
	return dashboarddto.OverviewDTO{
		MonitoredProjects:          value.MonitoredProjects,
		ActiveIncidents:            value.ActiveIncidents,
		WorkflowCompletedIncidents: value.WorkflowCompletedIncidents,
		PullRequestsCreated:        value.PullRequestsCreated,
		ActiveInvestigations:       incidents,
		LatestActivity:             activity,
	}
}

func ActivityEventToDTO(value domain.ActivityEvent) dashboarddto.ActivityEventDTO {
	var project *incidentdto.ProjectReferenceDTO
	if value.Project != nil {
		project = &incidentdto.ProjectReferenceDTO{ID: value.Project.ID, Name: value.Project.Name}
	}
	return dashboarddto.ActivityEventDTO{
		ID: value.ID, Type: string(value.Type), Project: project, IncidentID: value.IncidentID,
		IncidentKey: value.IncidentKey, OccurredAt: value.OccurredAt, Summary: value.Summary,
		UserMessage: value.UserMessage,
	}
}
