package domain

import (
	"time"

	"github.com/google/uuid"
)

type GitHubAppRegistrationStatus string

const (
	GitHubAppRegistrationCreated   GitHubAppRegistrationStatus = "created"
	GitHubAppRegistrationConverted GitHubAppRegistrationStatus = "converted"
	GitHubAppRegistrationCompleted GitHubAppRegistrationStatus = "completed"
)

type GitHubAppRegistration struct {
	ID                      uuid.UUID                   `gorm:"column:id;type:uuid;primaryKey"`
	DisplayName             string                      `gorm:"column:display_name"`
	AccountType             GitHubAccountType           `gorm:"column:account_type"`
	AccountIdentifier       string                      `gorm:"column:account_identifier"`
	ConversionStateDigest   []byte                      `gorm:"column:conversion_state_digest"`
	InstallationStateDigest []byte                      `gorm:"column:installation_state_digest"`
	Status                  GitHubAppRegistrationStatus `gorm:"column:status"`
	AppID                   *int64                      `gorm:"column:app_id"`
	AppSlug                 string                      `gorm:"column:app_slug"`
	AppName                 string                      `gorm:"column:app_name"`
	ClientID                string                      `gorm:"column:client_id"`
	ExpiresAt               time.Time                   `gorm:"column:expires_at"`
	ConversionConsumedAt    *time.Time                  `gorm:"column:conversion_consumed_at"`
	InstallationConsumedAt  *time.Time                  `gorm:"column:installation_consumed_at"`
	CreatedAt               time.Time                   `gorm:"column:created_at"`
	UpdatedAt               time.Time                   `gorm:"column:updated_at"`
}

type GitHubAppBinding struct {
	GitHubAccountID uuid.UUID `gorm:"column:github_account_id;type:uuid;primaryKey"`
	AppID           int64     `gorm:"column:app_id"`
	InstallationID  int64     `gorm:"column:installation_id"`
	AppSlug         string    `gorm:"column:app_slug"`
	ClientID        string    `gorm:"column:client_id"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}
