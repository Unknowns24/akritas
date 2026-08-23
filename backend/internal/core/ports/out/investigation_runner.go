package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
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
}

// InvestigationRunner performs the actual investigation work. The QVAC-backed
// implementation is out of scope for this task (pending PB-028+); production
// wiring uses a stub that always fails with an explicit "not implemented yet"
// message so the async pipeline (create -> queue -> poll) stays demonstrable.
type InvestigationRunner interface {
	Run(ctx context.Context, investigation domain.Investigation) (InvestigationRunResult, error)
}
