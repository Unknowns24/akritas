package evidence

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type UseCase struct {
	investigations portsout.InvestigationGetter
	evidence       portsout.EvidenceStore
}

func New(investigations portsout.InvestigationGetter, evidence portsout.EvidenceStore) portsin.EvidenceUseCase {
	return &UseCase{investigations: investigations, evidence: evidence}
}

func (uc *UseCase) ListInvestigationEvidence(ctx context.Context, investigationID uuid.UUID, params paging.Params) (paging.Slice[domain.Evidence], error) {
	if _, err := uc.investigations.FindByID(ctx, investigationID); err != nil {
		return paging.Slice[domain.Evidence]{}, err
	}
	return uc.evidence.ListByInvestigation(ctx, investigationID, params)
}
