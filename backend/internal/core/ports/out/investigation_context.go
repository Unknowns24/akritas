package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// InvestigationRunContext is the complete, sanitized application boundary
// handed to QVAC. External inference adapters must not query persistence.
type InvestigationRunContext struct {
	Investigation domain.Investigation
	Incident      domain.Incident
	Project       domain.Project
	Repository    RepositoryScope
	Evidence      []domain.Evidence
}

// InvestigationContextAssembler resolves H2 records and builds the initial
// bounded Evidence corpus before inference starts.
type InvestigationContextAssembler interface {
	Assemble(context.Context, domain.Investigation) (InvestigationRunContext, error)
}
