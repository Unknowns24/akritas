package investigation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) GetInvestigation(ctx context.Context, id uuid.UUID) (*domain.Investigation, error) {
	return uc.investigations.FindByID(ctx, id)
}
