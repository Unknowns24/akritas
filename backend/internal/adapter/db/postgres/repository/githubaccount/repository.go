package githubaccount

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	"gorm.io/gorm"
)

var ErrInvalidRepository = errors.New("invalid GitHub account repository configuration")

type Repository struct {
	db          *gorm.DB
	credentials *credentialstore.Store
}

func New(db *gorm.DB, credentials *credentialstore.Store) (*Repository, error) {
	if db == nil || credentials == nil {
		return nil, ErrInvalidRepository
	}
	return &Repository{db: db, credentials: credentials}, nil
}
