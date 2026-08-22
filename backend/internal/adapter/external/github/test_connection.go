package github

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) TestConnection(ctx context.Context, account domain.GitHubAccount) (portsout.ProviderConnectionResult, error) {
	startedAt := c.now()
	result := portsout.ProviderConnectionResult{CheckedAt: startedAt.UTC(), UserMessage: "La conexión con GitHub está disponible."}
	token, err := c.accountToken(ctx, account)
	if err != nil {
		result.Status = domain.ConnectionTestUnavailable
		result.UserMessage = "No se pudo acceder a la credencial de GitHub."
		return result, nil
	}
	defer wipe(token)
	path := "/user"
	if account.AuthenticationMethod == domain.GitHubAuthenticationGitHubApp {
		path = "/installation/repositories?per_page=1"
	} else if account.AccountType == domain.GitHubAccountOrganization {
		path = "/orgs/" + url.PathEscape(account.AccountIdentifier)
	}
	_, err = c.doJSON(ctx, http.MethodGet, path, string(token), nil, &struct{}{})
	result.Latency = c.now().Sub(startedAt)
	if err == nil {
		result.Status = domain.ConnectionTestConnected
		return result, nil
	}
	var providerErr *providerError
	if errors.As(err, &providerErr) && (providerErr.Status == http.StatusUnauthorized || providerErr.Status == http.StatusForbidden) {
		result.Status = domain.ConnectionTestAuthenticationFailed
		result.UserMessage = "GitHub rechazó la credencial configurada."
		return result, nil
	}
	result.Status = domain.ConnectionTestUnavailable
	result.UserMessage = "GitHub no está disponible temporalmente."
	return result, nil
}
