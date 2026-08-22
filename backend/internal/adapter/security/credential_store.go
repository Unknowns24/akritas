package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

var ErrEmptyMasterKey = errors.New("master key must not be empty")

// aesGCMCredentialStore implements the ADR-005 Credential Store. AKRITAS_MASTER_KEY
// is hashed into a 32-byte key so any high-entropy string the operator supplies works,
// regardless of its own length or encoding.
type aesGCMCredentialStore struct {
	gcm cipher.AEAD
}

func NewCredentialStore(masterKey string) (out.CredentialStore, error) {
	if masterKey == "" {
		return nil, ErrEmptyMasterKey
	}
	key := sha256.Sum256([]byte(masterKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesGCMCredentialStore{gcm: gcm}, nil
}

func (s *aesGCMCredentialStore) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return s.gcm.Seal(nonce, nonce, plaintext, nil), nil
}
