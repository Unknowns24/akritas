package mapper

import (
	investigationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/investigation"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// InvestigationToDTO always populates started_at: the contract marks it
// required, but a pending Investigation has no domain StartedAt yet, so this
// falls back to CreatedAt until Start() sets a real value.
func InvestigationToDTO(value domain.Investigation) investigationdto.InvestigationDTO {
	startedAt := value.CreatedAt
	if value.StartedAt != nil {
		startedAt = *value.StartedAt
	}
	dto := investigationdto.InvestigationDTO{
		ID: value.ID.String(), IncidentID: value.IncidentID.String(), Status: string(value.Status),
		StartedAt: startedAt, FinishedAt: value.FinishedAt, Summary: value.Summary, RootCause: value.RootCause,
		Confidence: value.Confidence, Hypotheses: value.Hypotheses, RelevantFiles: value.RelevantFiles,
		RelevantCommits: value.RelevantCommits, RecommendedActions: value.RecommendedActions,
		EvidenceCount: value.EvidenceCount, FailureUserMessage: value.FailureUserMessage,
	}
	if value.RootCauseStatus != nil {
		status := string(*value.RootCauseStatus)
		dto.RootCauseStatus = &status
	}
	if value.ResolutionStatus != nil {
		status := string(*value.ResolutionStatus)
		dto.ResolutionStatus = &status
	}
	return dto
}
