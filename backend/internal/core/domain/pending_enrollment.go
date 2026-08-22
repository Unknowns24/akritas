package domain

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type PendingEnrollment struct {
	ID                  uuid.UUID
	Email               string
	DisplayName         string
	PasswordHash        string
	EncryptedTOTPSecret []byte
	CreatedAt           time.Time
	ExpiresAt           time.Time
}

func NewPendingEnrollment(id uuid.UUID, email, displayName, passwordHash string, encryptedTOTPSecret []byte, createdAt, expiresAt time.Time) (*PendingEnrollment, error) {
	enrollment := &PendingEnrollment{
		ID: id, Email: email, DisplayName: displayName, PasswordHash: passwordHash,
		EncryptedTOTPSecret: encryptedTOTPSecret, CreatedAt: createdAt, ExpiresAt: expiresAt,
	}
	if err := enrollment.Validate(); err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (e PendingEnrollment) Validate() error {
	parsed, err := mail.ParseAddress(e.Email)
	invalid := e.ID == uuid.Nil || err != nil || parsed.Address != e.Email || !nonBlank(e.DisplayName) ||
		!nonBlank(e.PasswordHash) || len(e.EncryptedTOTPSecret) == 0 || !validTime(e.CreatedAt) || !e.ExpiresAt.After(e.CreatedAt)
	if invalid {
		return ErrInvalidPendingEnrollment.Wrap(validationCause("pending enrollment"))
	}
	return nil
}

func (e PendingEnrollment) IsExpired(now time.Time) bool {
	return !now.Before(e.ExpiresAt)
}
