package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

// EvidenceStore is write-once: Evidence is never updated after creation.
type EvidenceStore interface {
	Create(context.Context, *domain.Evidence) error
	ListByInvestigation(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.Evidence], error)
}
