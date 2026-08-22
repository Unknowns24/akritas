package github

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) VerifyInstallation(ctx context.Context, registration portsout.GitHubAppRegistration, installationID int64) (portsout.GitHubInstallation, error) {
	if installationID < 1 || registration.AppID == nil || *registration.AppID < 1 || c.credentials == nil {
		return portsout.GitHubInstallation{}, domain.ErrGitHubCredentialRejected
	}
	jwt, err := c.appJWT(ctx, portsout.CredentialOwnerGitHubManifest, registration.ID, *registration.AppID)
	if err != nil {
		return portsout.GitHubInstallation{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	var response struct {
		ID      int64 `json:"id"`
		Account struct {
			Login string `json:"login"`
			Type  string `json:"type"`
		} `json:"account"`
	}
	_, err = c.doJSON(ctx, http.MethodGet, "/app/installations/"+strconv.FormatInt(installationID, 10), jwt, nil, &response)
	if err != nil {
		return portsout.GitHubInstallation{}, normalizeCredentialError(err)
	}
	var accountType domain.GitHubAccountType
	switch {
	case strings.EqualFold(response.Account.Type, "User"):
		accountType = domain.GitHubAccountPersonal
	case strings.EqualFold(response.Account.Type, "Organization"):
		accountType = domain.GitHubAccountOrganization
	default:
		return portsout.GitHubInstallation{}, domain.ErrGitHubCredentialRejected
	}
	if response.ID != installationID || response.Account.Login == "" {
		return portsout.GitHubInstallation{}, domain.ErrGitHubCredentialRejected
	}
	return portsout.GitHubInstallation{InstallationID: response.ID, AccountLogin: response.Account.Login, AccountType: accountType}, nil
}
