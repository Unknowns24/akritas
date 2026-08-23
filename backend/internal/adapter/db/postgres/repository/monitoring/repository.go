package monitoring

import (
	"errors"

	"gorm.io/gorm"
)

var ErrInvalidRepository = errors.New("invalid monitoring repository configuration")

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) (*Repository, error) {
	if db == nil {
		return nil, ErrInvalidRepository
	}
	return &Repository{db: db}, nil
}
