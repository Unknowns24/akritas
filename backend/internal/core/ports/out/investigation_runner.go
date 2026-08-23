package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// InvestigationRunResult carries the structured classification an
// InvestigationRunner produced, ready for domain.Investigation.Complete.
type InvestigationRunResult struct {
	Summary            string
	RootCause          string
	RootCauseStatus    domain.RootCauseStatus
	ResolutionStatus   domain.ResolutionStatus
	Confidence         float64
	Hypotheses         []string
	RelevantFiles      []string
	RelevantCommits    []string
	RecommendedActions []string
	EvidenceIDs        []uuid.UUID
	DiscoveredEvidence []domain.Evidence
}

// InvestigationRunner performs local inference from an application-assembled
// context. Implementations must not query application persistence.
type InvestigationRunner interface {
	Run(ctx context.Context, runContext InvestigationRunContext) (InvestigationRunResult, error)
}
