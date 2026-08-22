package credentialstore

import (
	"context"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) DeleteOwnerTx(ctx context.Context, tx *gorm.DB, ownerType string, ownerID uuid.UUID) error {
	if err := tx.WithContext(ctx).Table("integration_credentials").Where("owner_type = ? AND owner_id = ?", ownerType, ownerID).Delete(&credentialRecord{}).Error; err != nil {
		return dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return nil
}
