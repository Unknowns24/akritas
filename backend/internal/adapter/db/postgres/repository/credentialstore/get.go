package credentialstore

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) Get(ctx context.Context, ownerType string, ownerID uuid.UUID, kind portsout.SecretKind) ([]byte, error) {
	var record credentialRecord
	err := s.db.WithContext(ctx).Table("integration_credentials").Where("owner_type = ? AND owner_id = ? AND secret_kind = ?", ownerType, ownerID, string(kind)).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrIntegrationNotFound
	}
	if err != nil {
		return nil, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	plaintext, err := s.cipher.Decrypt(credentialcipher.SealedValue{Ciphertext: record.Ciphertext, Nonce: record.Nonce}, credentialcipher.Metadata{
		OwnerType: record.OwnerType, OwnerID: record.OwnerID.String(), SecretKind: record.SecretKind, Version: record.EncryptionVersion,
	})
	if err != nil {
		return nil, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return plaintext, nil
}
