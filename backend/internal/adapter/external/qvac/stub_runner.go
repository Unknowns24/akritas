package qvac

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// ErrNotImplemented is retained for older tests that still assert against the stub.
var ErrNotImplemented = errors.New("QVAC integration is not implemented yet; see PB-028+")

// StubRunner always fails. Prefer Runner for production wiring.
type StubRunner struct{}

func NewStubRunner() *StubRunner { return &StubRunner{} }

func (r *StubRunner) Run(ctx context.Context, investigation domain.Investigation) (out.InvestigationRunResult, error) {
	return out.InvestigationRunResult{}, ErrNotImplemented
}
