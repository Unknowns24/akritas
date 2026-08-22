package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type StartGitHubAppRegistrationCommand struct {
	DisplayName  string
	OwnerType    domain.GitHubAccountType
	Organization string
}

type GitHubManifestRegistrationResult struct {
	RegistrationID uuid.UUID
	FormAction     string
	Manifest       string
	State          string
	ExpiresAt      string
}

type GitHubManifestCallbackResult struct {
	RedirectURL string
}

type GitHubInstallationCallbackResult struct {
	Account     domain.GitHubAccount
	RedirectURL string
}

type GitHubAppUseCase interface {
	StartRegistration(context.Context, StartGitHubAppRegistrationCommand) (GitHubManifestRegistrationResult, error)
	CompleteManifest(context.Context, string, string) (GitHubManifestCallbackResult, error)
	CompleteInstallation(context.Context, int64, string) (GitHubInstallationCallbackResult, error)
}
