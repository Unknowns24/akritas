package model

import (
	"time"

	"github.com/google/uuid"
)

type PendingEnrollment struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email               string    `gorm:"type:text;not null"`
	DisplayName         string    `gorm:"type:text;not null"`
	PasswordHash        string    `gorm:"type:text;not null"`
	EncryptedTOTPSecret []byte    `gorm:"type:bytea;not null"`
	CreatedAt           time.Time `gorm:"not null"`
	ExpiresAt           time.Time `gorm:"not null"`
}

func (PendingEnrollment) TableName() string {
	return "pending_enrollments"
}
