package model

import (
	"time"

	"github.com/google/uuid"
)

type Administrator struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"type:text;uniqueIndex;not null"`
	DisplayName  string    `gorm:"type:text;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (Administrator) TableName() string {
	return "administrators"
}
