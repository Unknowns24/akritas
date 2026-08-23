package out

import (
	"context"

	"github.com/google/uuid"
)

type SecretKind string

const (
	SecretKindGitHubPAT         SecretKind = "github_pat"
	SecretKindGitHubPrivateKey  SecretKind = "github_app_private_key"
	SecretKindGitHubWebhook     SecretKind = "github_app_webhook_secret"
	SecretKindDokployAPIKey     SecretKind = "dokploy_api_key"
	SecretKindAdministratorTOTP SecretKind = "administrator_totp"
	SecretKindQvacBearerToken   SecretKind = "qvac_bearer_token"
	SecretKindQvacBasicPassword SecretKind = "qvac_basic_password"
)

const (
	CredentialOwnerGitHubAccount     = "github_account"
	CredentialOwnerDokployServer     = "dokploy_server"
	CredentialOwnerGitHubManifest    = "github_app_registration"
	CredentialOwnerPendingEnrollment = "pending_enrollment"
	CredentialOwnerAdministrator     = "administrator"
	CredentialOwnerQvacConfiguration = "qvac_configuration"
)

type SecretValue struct {
	Kind      SecretKind
	Plaintext []byte
}

type CredentialStore interface {
	Put(context.Context, string, uuid.UUID, SecretValue) error
	Get(context.Context, string, uuid.UUID, SecretKind) ([]byte, error)
	DeleteOwner(context.Context, string, uuid.UUID) error
	MoveOwner(context.Context, string, uuid.UUID, string, uuid.UUID) error
}
