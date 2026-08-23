package mapper

import (
	investigationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/investigation"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func InvestigationToDTO(value domain.Investigation) investigationdto.InvestigationDTO {
	evidenceIDs := make([]string, 0, len(value.EvidenceIDs))
	for _, evidenceID := range value.EvidenceIDs {
		evidenceIDs = append(evidenceIDs, evidenceID.String())
	}
	dto := investigationdto.InvestigationDTO{
		ID: value.ID.String(), IncidentID: value.IncidentID.String(), Status: string(value.Status),
		CreatedAt: value.CreatedAt, StartedAt: value.StartedAt, FinishedAt: value.FinishedAt, Summary: value.Summary, RootCause: value.RootCause,
		Confidence: value.Confidence, Hypotheses: value.Hypotheses, RelevantFiles: value.RelevantFiles,
		RelevantCommits: value.RelevantCommits, RecommendedActions: value.RecommendedActions,
		EvidenceIDs: evidenceIDs, EvidenceCount: value.EvidenceCount, FailureUserMessage: value.FailureUserMessage,
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
