package credentialstore

import (
	"context"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func (s *Store) Get(ctx context.Context, ownerType string, ownerID uuid.UUID, kind portsout.SecretKind) ([]byte, error) {
	return s.GetWithDB(ctx, txcontext.From(ctx, s.db), ownerType, ownerID, kind)
}
