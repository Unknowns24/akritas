package credentialstore

import (
	"time"

	"github.com/google/uuid"
)

// credentialRecord is deliberately private infrastructure state. Ciphertext
// and nonce must never become domain fields or REST DTOs (ADR-005/ADR-012).
type credentialRecord struct {
	ID                uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	OwnerType         string    `gorm:"column:owner_type"`
	OwnerID           uuid.UUID `gorm:"column:owner_id;type:uuid"`
	SecretKind        string    `gorm:"column:secret_kind"`
	Ciphertext        []byte    `gorm:"column:ciphertext"`
	Nonce             []byte    `gorm:"column:nonce"`
	EncryptionVersion int       `gorm:"column:encryption_version"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}
