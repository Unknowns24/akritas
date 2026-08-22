package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type IntegrationStatus string

const (
	IntegrationStatusPending              IntegrationStatus = "pending"
	IntegrationStatusConnected            IntegrationStatus = "connected"
	IntegrationStatusAuthenticationFailed IntegrationStatus = "authentication_failed"
	IntegrationStatusUnavailable          IntegrationStatus = "unavailable"
)

func (s IntegrationStatus) Validate() error {
	switch s {
	case IntegrationStatusPending, IntegrationStatusConnected, IntegrationStatusAuthenticationFailed, IntegrationStatusUnavailable:
		return nil
	default:
		return ErrInvalidIntegrationStatus.Wrap(validationCause("integration status"))
	}
}

type ConnectionTestStatus string

const (
	ConnectionTestConnected            ConnectionTestStatus = "connected"
	ConnectionTestAuthenticationFailed ConnectionTestStatus = "authentication_failed"
	ConnectionTestUnavailable          ConnectionTestStatus = "unavailable"
)

func (s ConnectionTestStatus) Validate() error {
	switch s {
	case ConnectionTestConnected, ConnectionTestAuthenticationFailed, ConnectionTestUnavailable:
		return nil
	default:
		return ErrInvalidConnectionTestStatus.Wrap(validationCause("connection test status"))
	}
}

type GitHubAccountType string

const (
	GitHubAccountPersonal     GitHubAccountType = "personal"
	GitHubAccountOrganization GitHubAccountType = "organization"
)

func (t GitHubAccountType) Validate() error {
	switch t {
	case GitHubAccountPersonal, GitHubAccountOrganization:
		return nil
	default:
		return ErrInvalidGitHubAccount.Wrap(validationCause("GitHub account type"))
	}
}

type GitHubAuthenticationMethod string

const (
	GitHubAuthenticationPersonalAccessToken GitHubAuthenticationMethod = "personal_access_token"
	GitHubAuthenticationGitHubApp           GitHubAuthenticationMethod = "github_app"
)

func (m GitHubAuthenticationMethod) Validate() error {
	switch m {
	case GitHubAuthenticationPersonalAccessToken, GitHubAuthenticationGitHubApp:
		return nil
	default:
		return ErrInvalidGitHubAccount.Wrap(validationCause("GitHub authentication method"))
	}
}

type GitHubAccount struct {
	ID                   uuid.UUID                  `gorm:"column:id;type:uuid;primaryKey"`
	DisplayName          string                     `gorm:"column:display_name"`
	AccountType          GitHubAccountType          `gorm:"column:account_type"`
	AuthenticationMethod GitHubAuthenticationMethod `gorm:"column:authentication_method"`
	AccountIdentifier    string                     `gorm:"column:account_identifier"`
	AuthenticationStatus IntegrationStatus          `gorm:"column:authentication_status"`
	CredentialConfigured bool                       `gorm:"column:credential_configured"`
	RepositoryCount      int                        `gorm:"column:repository_count"`
	LastCheckedAt        *time.Time                 `gorm:"column:last_checked_at"`
	ManageURL            string                     `gorm:"column:manage_url"`
	CreatedAt            time.Time                  `gorm:"column:created_at"`
	UpdatedAt            time.Time                  `gorm:"column:updated_at"`
}

func NewGitHubAccount(
	id uuid.UUID,
	displayName string,
	accountType GitHubAccountType,
	authenticationMethod GitHubAuthenticationMethod,
	accountIdentifier string,
	authenticationStatus IntegrationStatus,
	createdAt time.Time,
) (*GitHubAccount, error) {
	account := &GitHubAccount{
		ID: id, DisplayName: strings.TrimSpace(displayName), AccountType: accountType,
		AuthenticationMethod: authenticationMethod, AccountIdentifier: strings.TrimSpace(accountIdentifier),
		AuthenticationStatus: authenticationStatus, CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	if err := account.Validate(); err != nil {
		return nil, err
	}
	return account, nil
}

func (a GitHubAccount) Validate() error {
	if a.ID == uuid.Nil || !nonBlank(a.DisplayName) || !nonBlank(a.AccountIdentifier) ||
		a.AccountType.Validate() != nil || a.AuthenticationMethod.Validate() != nil ||
		a.AuthenticationStatus.Validate() != nil || a.RepositoryCount < 0 || !validTime(a.CreatedAt) || a.UpdatedAt.Before(a.CreatedAt) {
		return ErrInvalidGitHubAccount.Wrap(validationCause("GitHub account"))
	}
	if a.ManageURL != "" && !validHTTPURL(a.ManageURL) {
		return ErrInvalidGitHubAccount.Wrap(validationCause("GitHub manage URL"))
	}
	return nil
}
