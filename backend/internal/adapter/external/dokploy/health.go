package dokploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (c *Client) health(ctx context.Context, rawURL, credential string) (domain.ConnectionTestStatus, error) {
	base, err := c.normalizeAndValidateURL(ctx, rawURL)
	if err != nil {
		return domain.ConnectionTestUnavailable, domain.ErrIntegrationUnavailable
	}
	_, err = c.do(ctx, base+"/api/settings.health", credential)
	if err == nil {
		return domain.ConnectionTestConnected, nil
	}
	var providerErr *providerError
	if errors.As(err, &providerErr) && (providerErr.Status == http.StatusUnauthorized || providerErr.Status == http.StatusForbidden) {
		return domain.ConnectionTestAuthenticationFailed, nil
	}
	return domain.ConnectionTestUnavailable, domain.ErrIntegrationUnavailable.Wrap(err)
}
