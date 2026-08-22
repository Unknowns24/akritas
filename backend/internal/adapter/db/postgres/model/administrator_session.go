package model

import (
	"time"

	"github.com/google/uuid"
)

type AdministratorSession struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	AdministratorID   uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash         string    `gorm:"type:text;not null;uniqueIndex"`
	AuthenticatedAt   time.Time `gorm:"not null"`
	IdleExpiresAt     time.Time `gorm:"not null"`
	AbsoluteExpiresAt time.Time `gorm:"not null"`
	RevokedAt         *time.Time
}

func (AdministratorSession) TableName() string {
	return "administrator_sessions"
}
