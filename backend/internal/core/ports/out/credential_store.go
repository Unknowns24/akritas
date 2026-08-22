package out

import (
	"context"

	"github.com/google/uuid"
)

type SecretKind string

const (
	SecretKindGitHubPAT        SecretKind = "github_pat"
	SecretKindGitHubPrivateKey SecretKind = "github_app_private_key"
	SecretKindGitHubWebhook    SecretKind = "github_app_webhook_secret"
	SecretKindDokployAPIKey    SecretKind = "dokploy_api_key"
)

const (
	CredentialOwnerGitHubAccount  = "github_account"
	CredentialOwnerDokployServer  = "dokploy_server"
	CredentialOwnerGitHubManifest = "github_app_registration"
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
