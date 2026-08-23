package credentialstore

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Store) PutTx(ctx context.Context, tx *gorm.DB, ownerType string, ownerID uuid.UUID, secret portsout.SecretValue) error {
	metadata := credentialcipher.Metadata{OwnerType: ownerType, OwnerID: ownerID.String(), SecretKind: string(secret.Kind), Version: encryptionVersion}
	sealed, err := s.cipher.Encrypt(secret.Plaintext, metadata)
	if err != nil {
		return dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	now := time.Now().UTC()
	record := credentialRecord{ID: uuid.New(), OwnerType: ownerType, OwnerID: ownerID, SecretKind: string(secret.Kind), Ciphertext: sealed.Ciphertext, Nonce: sealed.Nonce, EncryptionVersion: encryptionVersion, CreatedAt: now, UpdatedAt: now}
	err = tx.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "owner_type"}, {Name: "owner_id"}, {Name: "secret_kind"}},
		DoUpdates: clause.Assignments(map[string]any{"ciphertext": sealed.Ciphertext, "nonce": sealed.Nonce, "encryption_version": encryptionVersion, "updated_at": now}),
	}).Table("credentials").Create(&record).Error
	if err != nil {
		return dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return nil
}
