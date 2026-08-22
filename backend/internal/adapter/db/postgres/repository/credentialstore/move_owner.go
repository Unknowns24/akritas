package credentialstore

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) MoveOwner(ctx context.Context, fromType string, fromID uuid.UUID, toType string, toID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.MoveOwnerTx(ctx, tx, fromType, fromID, toType, toID)
	})
}
