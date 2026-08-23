package mapper

import (
	evidencedto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/evidence"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func EvidenceToDTO(value domain.Evidence) evidencedto.EvidenceDTO {
	return evidencedto.EvidenceDTO{
		ID: value.ID.String(), InvestigationID: value.InvestigationID.String(), Type: string(value.Type),
		Summary: value.Summary, Content: value.Content, FilePath: value.FilePath,
		LineStart: value.LineStart, LineEnd: value.LineEnd, CommitSHA: value.CommitSHA, Patch: value.Patch,
		OccurredAt: value.OccurredAt, Redacted: value.Redacted, CreatedAt: value.CreatedAt,
	}
}
