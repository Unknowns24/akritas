package credentialstore

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) MoveOwner(ctx context.Context, fromType string, fromID uuid.UUID, toType string, toID uuid.UUID) error {
	return txcontext.From(ctx, s.db).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.MoveOwnerTx(ctx, tx, fromType, fromID, toType, toID)
	})
}
