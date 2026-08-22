package domain

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Administrator struct {
	ID          uuid.UUID
	Email       string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAdministrator(id uuid.UUID, email, displayName string, createdAt time.Time) (*Administrator, error) {
	administrator := &Administrator{
		ID: id, Email: strings.TrimSpace(email), DisplayName: strings.TrimSpace(displayName),
		CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	if err := administrator.Validate(); err != nil {
		return nil, err
	}
	return administrator, nil
}

func (a Administrator) Validate() error {
	parsed, err := mail.ParseAddress(a.Email)
	if a.ID == uuid.Nil || err != nil || parsed.Address != a.Email || !nonBlank(a.DisplayName) || !validTime(a.CreatedAt) || a.UpdatedAt.Before(a.CreatedAt) {
		return ErrInvalidAdministrator.Wrap(validationCause("administrator"))
	}
	return nil
}
