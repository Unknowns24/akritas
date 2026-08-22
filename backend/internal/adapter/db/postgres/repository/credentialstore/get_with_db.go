package credentialstore

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) GetWithDB(ctx context.Context, tx *gorm.DB, ownerType string, ownerID uuid.UUID, kind portsout.SecretKind) ([]byte, error) {
	var record credentialRecord
	if err := tx.WithContext(ctx).Table("credentials").Where("owner_type = ? AND owner_id = ? AND secret_kind = ?", ownerType, ownerID, string(kind)).First(&record).Error; err != nil {
		return nil, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	sealed := credentialcipher.SealedValue{Ciphertext: record.Ciphertext, Nonce: record.Nonce}
	metadata := credentialcipher.Metadata{OwnerType: record.OwnerType, OwnerID: record.OwnerID.String(), SecretKind: record.SecretKind, Version: record.EncryptionVersion}
	plaintext, err := s.cipher.Decrypt(sealed, metadata)
	if err != nil {
		return nil, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return plaintext, nil
}
