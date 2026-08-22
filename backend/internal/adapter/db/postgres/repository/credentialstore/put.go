package credentialstore

import (
	"context"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

const encryptionVersion = 1

func (s *Store) Put(ctx context.Context, ownerType string, ownerID uuid.UUID, secret portsout.SecretValue) error {
	return s.PutTx(ctx, s.db, ownerType, ownerID, secret)
}
