package githubapp

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (uc *UseCase) CompleteInstallation(ctx context.Context, installationID int64, state string) (portsin.GitHubInstallationCallbackResult, error) {
	if installationID < 1 || len(state) < 32 || len(state) > 512 {
		return portsin.GitHubInstallationCallbackResult{}, domain.ErrManifestStateInvalid
	}
	digest := sha256.Sum256([]byte(state))
	registration, err := uc.store.ConsumeInstallationState(ctx, digest[:], uc.now().UTC())
	if err != nil {
		return portsin.GitHubInstallationCallbackResult{}, err
	}
	if registration.AppID == nil {
		return portsin.GitHubInstallationCallbackResult{}, domain.ErrManifestStateInvalid
	}
	installation, err := uc.gateway.VerifyInstallation(ctx, *registration, installationID)
	if err != nil {
		return portsin.GitHubInstallationCallbackResult{}, err
	}
	if installation.InstallationID != installationID || installation.AccountType != registration.AccountType || (registration.AccountType == domain.GitHubAccountOrganization && !strings.EqualFold(installation.AccountLogin, registration.AccountIdentifier)) {
		return portsin.GitHubInstallationCallbackResult{}, domain.ErrGitHubCredentialRejected
	}
	now := uc.now().UTC()
	account, err := domain.NewGitHubAccount(uc.newID(), registration.DisplayName, installation.AccountType, domain.GitHubAuthenticationGitHubApp, installation.AccountLogin, domain.IntegrationStatusConnected, now)
	if err != nil {
		return portsin.GitHubInstallationCallbackResult{}, err
	}
	account.CredentialConfigured = true
	account.LastCheckedAt = &now
	account.ManageURL = fmt.Sprintf("https://github.com/settings/installations/%d", installationID)
	binding := portsout.GitHubAppBinding{GitHubAccountID: account.ID, AppID: *registration.AppID, InstallationID: installationID, AppSlug: registration.AppSlug, ClientID: registration.ClientID}
	registration.Status = portsout.GitHubAppRegistrationCompleted
	registration.UpdatedAt = now
	if err := uc.store.CompleteInstallation(ctx, registration, account, binding); err != nil {
		return portsin.GitHubInstallationCallbackResult{}, err
	}
	redirect := strings.TrimRight(uc.publicURL.String(), "/") + "/settings/integrations/github?connected=" + account.ID.String()
	return portsin.GitHubInstallationCallbackResult{Account: *account, RedirectURL: redirect}, nil
}
