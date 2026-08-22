package credentialcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const keySize = 32

var (
	ErrInvalidKey       = errors.New("invalid credential encryption key")
	ErrInvalidMetadata  = errors.New("invalid credential encryption metadata")
	ErrEncryptionFailed = errors.New("credential encryption failed")
	ErrDecryptionFailed = errors.New("credential decryption failed")
)

type Metadata struct {
	OwnerType  string `json:"owner_type"`
	OwnerID    string `json:"owner_id"`
	SecretKind string `json:"secret_kind"`
	Version    int    `json:"version"`
}

func (m Metadata) validate() error {
	if m.OwnerType == "" || m.OwnerID == "" || m.SecretKind == "" || m.Version < 1 {
		return ErrInvalidMetadata
	}
	return nil
}

type SealedValue struct {
	Ciphertext []byte
	Nonce      []byte
}

type Cipher struct {
	aead cipher.AEAD
}

func New(encodedKey string) (*Cipher, error) {
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil || len(key) != keySize {
		return nil, ErrInvalidKey
	}
	return NewFromKey(key)
}

func NewFromKey(key []byte) (*Cipher, error) {
	if len(key) != keySize {
		return nil, ErrInvalidKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, ErrInvalidKey
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrInvalidKey
	}
	return &Cipher{aead: aead}, nil
}

func (c *Cipher) Encrypt(plaintext []byte, metadata Metadata) (SealedValue, error) {
	if c == nil || c.aead == nil || len(plaintext) == 0 {
		return SealedValue{}, ErrEncryptionFailed
	}
	aad, err := additionalData(metadata)
	if err != nil {
		return SealedValue{}, err
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return SealedValue{}, ErrEncryptionFailed
	}
	ciphertext := c.aead.Seal(nil, nonce, plaintext, aad)
	return SealedValue{Ciphertext: ciphertext, Nonce: nonce}, nil
}

func (c *Cipher) Decrypt(sealed SealedValue, metadata Metadata) ([]byte, error) {
	if c == nil || c.aead == nil || len(sealed.Ciphertext) == 0 || len(sealed.Nonce) != c.aead.NonceSize() {
		return nil, ErrDecryptionFailed
	}
	aad, err := additionalData(metadata)
	if err != nil {
		return nil, err
	}
	plaintext, err := c.aead.Open(nil, sealed.Nonce, sealed.Ciphertext, aad)
	if err != nil {
		return nil, ErrDecryptionFailed
	}
	return plaintext, nil
}

func additionalData(metadata Metadata) ([]byte, error) {
	if err := metadata.validate(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrInvalidMetadata)
	}
	return data, nil
}
