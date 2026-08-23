package mapper

import (
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestInvestigationToDTOKeepsPendingStartedAtOptionalAndExposesEvidenceIDs(t *testing.T) {
	t.Parallel()
	evidenceID := uuid.New()
	value := domain.Investigation{
		ID: uuid.New(), IncidentID: uuid.New(), Status: domain.InvestigationStatusPending,
		CreatedAt: time.Now().UTC(), EvidenceIDs: []uuid.UUID{evidenceID}, EvidenceCount: 1,
		Hypotheses: []string{}, RelevantFiles: []string{}, RelevantCommits: []string{}, RecommendedActions: []string{},
	}
	dto := InvestigationToDTO(value)
	if dto.CreatedAt != value.CreatedAt || dto.StartedAt != nil || len(dto.EvidenceIDs) != 1 || dto.EvidenceIDs[0] != evidenceID.String() {
		t.Fatalf("dto=%+v", dto)
	}
}
