package mapper

import (
	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	remediationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/remediation"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func RemediationToDTO(value domain.Remediation) remediationdto.RemediationDTO {
	changes := make([]remediationdto.CodeChangeDTO, 0, len(value.Changes))
	for _, change := range value.Changes {
		changes = append(changes, remediationdto.CodeChangeDTO{
			FilePath: string(change.FilePath), ChangeType: string(change.ChangeType), Patch: change.Patch, Redacted: true,
		})
	}
	summary := remediationdto.ValidationSummaryDTO{}
	for _, result := range value.ValidationResults {
		summary.Total++
		if result.Status == domain.ValidationStatusPassed {
			summary.Passed++
		}
		if result.Status == domain.ValidationStatusFailed {
			summary.Failed++
		}
	}
	var pr *remediationdto.PullRequestReferenceDTO
	if value.PullRequestReference != nil {
		pr = PullRequestReferenceToDTO(*value.PullRequestReference)
	}
	return remediationdto.RemediationDTO{
		ID: value.ID, IncidentID: value.IncidentID, Status: string(value.Status), BranchName: value.BranchName,
		ChangesSummary: value.ChangesSummary, Changes: changes, ValidationSummary: summary,
		FailureUserMessage: value.FailureUserMessage, PullRequestReference: pr,
		CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}
}

func ValidationResultToDTO(value domain.ValidationResult) remediationdto.ValidationResultDTO {
	return remediationdto.ValidationResultDTO{
		ID: value.ID, RemediationID: value.RemediationID, Type: string(value.Type), Name: value.Name,
		Status: string(value.Status), StartedAt: value.StartedAt, FinishedAt: value.FinishedAt,
		Summary: value.Summary, OutputExcerpt: value.OutputExcerpt, OutputRedacted: value.OutputRedacted,
	}
}

func PullRequestReferenceToDTO(value domain.PullRequestReference) *remediationdto.PullRequestReferenceDTO {
	return &remediationdto.PullRequestReferenceDTO{
		Number: value.Number, URL: value.URL, Repository: value.Repository, Branch: value.Branch, CreatedAt: value.CreatedAt,
	}
}

func PullRequestProjectionToDTO(value domain.PullRequestProjection) remediationdto.PullRequestDTO {
	var issue *incidentdto.GitHubIssueReferenceDTO
	if value.IssueReference != nil {
		issue = &incidentdto.GitHubIssueReferenceDTO{
			Number: value.IssueReference.Number, URL: value.IssueReference.URL,
			Repository: value.IssueReference.Repository, CreatedAt: value.IssueReference.CreatedAt,
		}
	}
	return remediationdto.PullRequestDTO{
		ID: value.ID, Project: incidentdto.ProjectReferenceDTO{ID: value.Project.ID, Name: value.Project.Name},
		IncidentID: value.IncidentID, IncidentKey: value.IncidentKey, RemediationID: value.RemediationID,
		IssueReference: issue, Reference: *PullRequestReferenceToDTO(value.Reference), Title: value.Title,
		ChangesSummary: value.ChangesSummary, ValidationSummary: value.ValidationSummary,
		RiskSummary: value.RiskSummary, CreatedAt: value.CreatedAt,
	}
}
