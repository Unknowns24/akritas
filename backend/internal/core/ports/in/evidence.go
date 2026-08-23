package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type EvidenceUseCase interface {
	ListInvestigationEvidence(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.Evidence], error)
}
