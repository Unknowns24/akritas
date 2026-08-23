package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	if id == uuid.Nil {
		return nil, domain.ErrIncidentNotFound
	}
	return uc.incidents.Get(ctx, id)
}
