package credentialstore

import (
	"context"

	"github.com/google/uuid"
)

func (s *Store) DeleteOwner(ctx context.Context, ownerType string, ownerID uuid.UUID) error {
	return s.DeleteOwnerTx(ctx, s.db, ownerType, ownerID)
}
