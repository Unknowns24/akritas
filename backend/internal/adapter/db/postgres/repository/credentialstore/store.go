package credentialstore

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	"gorm.io/gorm"
)

var ErrInvalidCredentialStore = errors.New("invalid credential store configuration")

type Store struct {
	db     *gorm.DB
	cipher *credentialcipher.Cipher
}

func New(db *gorm.DB, cipher *credentialcipher.Cipher) (*Store, error) {
	if db == nil || cipher == nil {
		return nil, ErrInvalidCredentialStore
	}
	return &Store{db: db, cipher: cipher}, nil
}
