package domain

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PendingEnrollment struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Email       string    `gorm:"column:email"`
	DisplayName string    `gorm:"column:display_name"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	ExpiresAt   time.Time `gorm:"column:expires_at"`
}

func NewPendingEnrollment(id uuid.UUID, email, displayName string, createdAt, expiresAt time.Time) (*PendingEnrollment, error) {
	enrollment := &PendingEnrollment{
		ID: id, Email: strings.TrimSpace(email), DisplayName: strings.TrimSpace(displayName),
		CreatedAt: createdAt, ExpiresAt: expiresAt,
	}
	if err := enrollment.Validate(); err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (e PendingEnrollment) Validate() error {
	parsed, err := mail.ParseAddress(e.Email)
	invalid := e.ID == uuid.Nil || err != nil || parsed.Address != e.Email || !nonBlank(e.DisplayName) ||
		!validTime(e.CreatedAt) || !e.ExpiresAt.After(e.CreatedAt)
	if invalid {
		return ErrInvalidPendingEnrollment.Wrap(validationCause("pending enrollment"))
	}
	return nil
}

func (e PendingEnrollment) IsExpired(now time.Time) bool {
	return !now.Before(e.ExpiresAt)
}
