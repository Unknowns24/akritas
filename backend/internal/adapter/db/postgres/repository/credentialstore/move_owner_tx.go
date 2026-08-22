package credentialstore

import (
	"context"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) MoveOwnerTx(ctx context.Context, tx *gorm.DB, fromType string, fromID uuid.UUID, toType string, toID uuid.UUID) error {
	var records []credentialRecord
	if err := tx.WithContext(ctx).Table("integration_credentials").Where("owner_type = ? AND owner_id = ?", fromType, fromID).Find(&records).Error; err != nil {
		return dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	for _, record := range records {
		plaintext, err := s.GetWithDB(ctx, tx, fromType, fromID, portsout.SecretKind(record.SecretKind))
		if err != nil {
			return err
		}
		if err := s.PutTx(ctx, tx, toType, toID, portsout.SecretValue{Kind: portsout.SecretKind(record.SecretKind), Plaintext: plaintext}); err != nil {
			wipe(plaintext)
			return err
		}
		wipe(plaintext)
	}
	return s.DeleteOwnerTx(ctx, tx, fromType, fromID)
}
