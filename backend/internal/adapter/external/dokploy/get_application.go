package dokploy

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) GetApplication(ctx context.Context, server domain.DokployServer, identifier string) (domain.DokployApplication, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" || c.credentials == nil {
		return domain.DokployApplication{}, domain.ErrIntegrationNotFound
	}
	credential, err := c.credentials.Get(ctx, portsout.CredentialOwnerDokployServer, server.ID, portsout.SecretKindDokployAPIKey)
	if err != nil {
		return domain.DokployApplication{}, domain.ErrIntegrationUnavailable
	}
	defer wipe(credential)
	base, err := c.normalizeAndValidateURL(ctx, server.BaseURL)
	if err != nil {
		return domain.DokployApplication{}, domain.ErrIntegrationUnavailable
	}
	query := url.Values{"q": {identifier}, "limit": {"100"}, "offset": {"0"}}
	response, err := c.do(ctx, base+"/api/application.search?"+query.Encode(), string(credential))
	if err != nil {
		return domain.DokployApplication{}, normalizeProviderError(err)
	}
	var payload applicationSearchResponse
	if err := json.Unmarshal(response, &payload); err != nil {
		if directErr := json.Unmarshal(response, &payload.Items); directErr != nil {
			return domain.DokployApplication{}, domain.ErrIntegrationUnavailable
		}
	}
	for _, item := range payload.Items {
		if item.ApplicationID != identifier {
			continue
		}
		return domain.NewDokployApplication(server.ID, item.ApplicationID, item.AppName, item.displayName(), item.environmentName(), mapStatus(item.status()))
	}
	return domain.DokployApplication{}, domain.ErrIntegrationNotFound
}
