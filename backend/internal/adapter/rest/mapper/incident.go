package mapper

import (
	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func IncidentToDTO(value domain.Incident) incidentdto.IncidentDTO {
	project := incidentdto.ProjectReferenceDTO{ID: value.ProjectID}
	if value.Project != nil {
		project = incidentdto.ProjectReferenceDTO{ID: value.Project.ID, Name: value.Project.Name}
	}
	result := incidentdto.IncidentDTO{
		ID: value.ID, Key: value.Key, Project: project, Fingerprint: value.Fingerprint,
		Severity: value.Severity, Title: value.Title, Summary: value.Summary, Phase: value.Phase,
		TerminalOutcome: value.TerminalOutcome, RootCauseStatus: value.RootCauseStatus,
		ResolutionStatus: value.ResolutionStatus, Confidence: value.Confidence,
		OccurrenceCount: value.OccurrenceCount, FirstSeenAt: value.FirstSeenAt, LastSeenAt: value.LastSeenAt,
	}
	if value.GitHubIssueReference != nil {
		result.GitHubIssueReference = githubIssueReferenceToDTO(*value.GitHubIssueReference)
	}
	if value.LatestInvestigation != nil {
		latest := InvestigationToDTO(*value.LatestInvestigation)
		result.LatestInvestigation = &latest
	}
	return result
}

func githubIssueReferenceToDTO(value domain.GitHubIssueReference) *incidentdto.GitHubIssueReferenceDTO {
	return &incidentdto.GitHubIssueReferenceDTO{
		Number: value.Number, URL: value.URL, Repository: value.Repository, CreatedAt: value.CreatedAt,
	}
}
