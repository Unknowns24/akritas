package github

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) accountToken(ctx context.Context, account domain.GitHubAccount) ([]byte, error) {
	if c.credentials == nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	if account.AuthenticationMethod == domain.GitHubAuthenticationPersonalAccessToken {
		return c.credentials.Get(ctx, portsout.CredentialOwnerGitHubAccount, account.ID, portsout.SecretKindGitHubPAT)
	}
	if account.AuthenticationMethod == domain.GitHubAuthenticationGitHubApp {
		return c.installationToken(ctx, account.ID)
	}
	return nil, domain.ErrIntegrationUnavailable
}
