package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type DokployServerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.DokployServer, error)
}
