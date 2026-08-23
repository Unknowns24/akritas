// Package qvac will host the local-only QVAC inference client (ADR-001).
// Real investigation execution is out of scope for AKR-INVESTIGATION-LIFECYCLE
// (pending PB-028+); StubRunner keeps the async pipeline demonstrable in the
// meantime by always failing with an explicit, non-technical message.
package qvac

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

var ErrNotImplemented = errors.New("QVAC integration is not implemented yet; see PB-028+")

type StubRunner struct{}

func NewStubRunner() *StubRunner {
	return &StubRunner{}
}

func (r *StubRunner) Run(ctx context.Context, investigation domain.Investigation) (out.InvestigationRunResult, error) {
	return out.InvestigationRunResult{}, ErrNotImplemented
}
