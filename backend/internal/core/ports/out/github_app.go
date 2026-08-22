package out

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type GitHubAppRegistrationStatus = domain.GitHubAppRegistrationStatus

const (
	GitHubAppRegistrationCreated   = domain.GitHubAppRegistrationCreated
	GitHubAppRegistrationConverted = domain.GitHubAppRegistrationConverted
	GitHubAppRegistrationCompleted = domain.GitHubAppRegistrationCompleted
)

type GitHubAppRegistration = domain.GitHubAppRegistration
type GitHubAppBinding = domain.GitHubAppBinding

type GitHubAppRegistrationStore interface {
	CreateRegistration(context.Context, *GitHubAppRegistration) error
	ConsumeConversionState(context.Context, []byte, time.Time) (*GitHubAppRegistration, error)
	CompleteConversion(context.Context, *GitHubAppRegistration, []SecretValue) error
	ConsumeInstallationState(context.Context, []byte, time.Time) (*GitHubAppRegistration, error)
	CompleteInstallation(context.Context, *GitHubAppRegistration, *domain.GitHubAccount, GitHubAppBinding) error
}

type GitHubAppBindingReader interface {
	GetBinding(context.Context, uuid.UUID) (GitHubAppBinding, error)
}

type GitHubManifestConversion struct {
	AppID         int64
	AppSlug       string
	AppName       string
	ClientID      string
	PrivateKey    []byte
	WebhookSecret []byte
}

type GitHubInstallation struct {
	InstallationID int64
	AccountLogin   string
	AccountType    domain.GitHubAccountType
}

type GitHubAppGateway interface {
	ExchangeManifest(context.Context, string) (GitHubManifestConversion, error)
	VerifyInstallation(context.Context, GitHubAppRegistration, int64) (GitHubInstallation, error)
}
