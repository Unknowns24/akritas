package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// EvidenceAssembler gathers whatever Evidence can be produced for an
// Investigation from sources that are actually available today. It returns
// an empty slice (not an error) when a source is simply absent (e.g. the
// incident or project is not found); a non-nil error means a real
// infrastructure failure.
type EvidenceAssembler interface {
	Assemble(ctx context.Context, investigation domain.Investigation) ([]domain.Evidence, error)
}
