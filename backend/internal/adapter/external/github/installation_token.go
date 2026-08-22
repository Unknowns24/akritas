package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func (c *Client) installationToken(ctx context.Context, accountID uuid.UUID) ([]byte, error) {
	if c.bindings == nil || c.credentials == nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	cacheKey := accountID.String()
	c.cacheMu.Lock()
	cached, ok := c.tokenCache[cacheKey]
	if ok && c.now().Add(time.Minute).Before(cached.expiresAt) {
		token := []byte(cached.token)
		c.cacheMu.Unlock()
		return token, nil
	}
	delete(c.tokenCache, cacheKey)
	c.cacheMu.Unlock()
	binding, err := c.bindings.GetBinding(ctx, accountID)
	if err != nil {
		return nil, err
	}
	jwt, err := c.appJWT(ctx, portsout.CredentialOwnerGitHubAccount, accountID, binding.AppID)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	var response struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	_, err = c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/app/installations/%d/access_tokens", binding.InstallationID), jwt, nil, &response)
	if err != nil {
		return nil, normalizeProviderError(err)
	}
	if response.Token == "" || !response.ExpiresAt.After(c.now()) {
		return nil, domain.ErrIntegrationUnavailable
	}
	c.cacheMu.Lock()
	c.tokenCache[cacheKey] = cachedInstallationToken{token: response.Token, expiresAt: response.ExpiresAt}
	c.cacheMu.Unlock()
	return []byte(response.Token), nil
}
