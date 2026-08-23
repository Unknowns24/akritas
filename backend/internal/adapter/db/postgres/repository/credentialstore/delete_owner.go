package credentialstore

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/google/uuid"
)

func (s *Store) DeleteOwner(ctx context.Context, ownerType string, ownerID uuid.UUID) error {
	return s.DeleteOwnerTx(ctx, txcontext.From(ctx, s.db), ownerType, ownerID)
}
