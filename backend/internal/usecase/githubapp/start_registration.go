package githubapp

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type manifestDocument struct {
	Name               string            `json:"name"`
	URL                string            `json:"url"`
	RedirectURL        string            `json:"redirect_url"`
	SetupURL           string            `json:"setup_url"`
	Public             bool              `json:"public"`
	HookAttributes     map[string]bool   `json:"hook_attributes"`
	DefaultPermissions map[string]string `json:"default_permissions"`
}

func (uc *UseCase) StartRegistration(ctx context.Context, command portsin.StartGitHubAppRegistrationCommand) (portsin.GitHubManifestRegistrationResult, error) {
	displayName := strings.TrimSpace(command.DisplayName)
	organization := strings.TrimSpace(command.Organization)
	if displayName == "" || len(displayName) > 100 || len(organization) > 100 || command.OwnerType.Validate() != nil || (command.OwnerType == domain.GitHubAccountOrganization && organization == "") || (command.OwnerType == domain.GitHubAccountPersonal && organization != "") {
		return portsin.GitHubManifestRegistrationResult{}, domain.ErrInvalidGitHubAccount
	}
	state, err := uc.newState()
	if err != nil {
		return portsin.GitHubManifestRegistrationResult{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	digest := sha256.Sum256([]byte(state))
	now := uc.now().UTC()
	registration := portsout.GitHubAppRegistration{
		ID: uc.newID(), DisplayName: displayName, AccountType: command.OwnerType, AccountIdentifier: organization,
		ConversionStateDigest: digest[:], Status: portsout.GitHubAppRegistrationCreated,
		ExpiresAt: now.Add(time.Hour), CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.store.CreateRegistration(ctx, &registration); err != nil {
		return portsin.GitHubManifestRegistrationResult{}, err
	}
	base := strings.TrimRight(uc.publicURL.String(), "/")
	manifest := manifestDocument{
		Name: fmt.Sprintf("Akritas-%s", registration.ID.String()[:12]), URL: base,
		RedirectURL: base + "/api/v1/integrations/github/app-manifest/callback",
		SetupURL:    base + "/api/v1/integrations/github/app-installations/callback",
		Public:      false, HookAttributes: map[string]bool{"active": false},
		DefaultPermissions: map[string]string{"metadata": "read", "contents": "write", "issues": "write", "pull_requests": "write"},
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return portsin.GitHubManifestRegistrationResult{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	baseFormAction := "https://github.com/settings/apps/new"
	if command.OwnerType == domain.GitHubAccountOrganization {
		baseFormAction = "https://github.com/organizations/" + url.PathEscape(organization) + "/settings/apps/new"
	}
	formURL, err := url.Parse(baseFormAction)
	if err != nil {
		return portsin.GitHubManifestRegistrationResult{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	query := formURL.Query()
	query.Set("state", state)
	formURL.RawQuery = query.Encode()
	return portsin.GitHubManifestRegistrationResult{
		RegistrationID: registration.ID, FormAction: formURL.String(), Manifest: string(manifestJSON), State: state, ExpiresAt: registration.ExpiresAt.Format(time.RFC3339),
	}, nil
}
