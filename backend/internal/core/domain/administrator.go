package domain

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Administrator struct {
	ID                     uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Email                  string    `gorm:"column:email"`
	DisplayName            string    `gorm:"column:display_name"`
	LastAcceptedTOTPPeriod int64     `gorm:"column:last_accepted_totp_period"`
	CreatedAt              time.Time `gorm:"column:created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at"`
}

func NewAdministrator(id uuid.UUID, email, displayName string, createdAt time.Time) (*Administrator, error) {
	administrator := &Administrator{
		ID: id, Email: strings.TrimSpace(email), DisplayName: strings.TrimSpace(displayName),
		LastAcceptedTOTPPeriod: -1, CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	if err := administrator.Validate(); err != nil {
		return nil, err
	}
	return administrator, nil
}

func (a Administrator) Validate() error {
	parsed, err := mail.ParseAddress(a.Email)
	if a.ID == uuid.Nil || err != nil || parsed.Address != a.Email || !nonBlank(a.DisplayName) || a.LastAcceptedTOTPPeriod < -1 || !validTime(a.CreatedAt) || a.UpdatedAt.Before(a.CreatedAt) {
		return ErrInvalidAdministrator.Wrap(validationCause("administrator"))
	}
	return nil
}
