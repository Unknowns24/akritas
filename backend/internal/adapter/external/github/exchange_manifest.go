package github

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) ExchangeManifest(ctx context.Context, code string) (portsout.GitHubManifestConversion, error) {
	if len(code) < 20 || len(code) > 255 {
		return portsout.GitHubManifestConversion{}, domain.ErrManifestStateInvalid
	}
	var response struct {
		ID            int64  `json:"id"`
		Slug          string `json:"slug"`
		Name          string `json:"name"`
		ClientID      string `json:"client_id"`
		PEM           string `json:"pem"`
		WebhookSecret string `json:"webhook_secret"`
	}
	_, err := c.doJSON(ctx, http.MethodPost, "/app-manifests/"+url.PathEscape(code)+"/conversions", "", nil, &response)
	if err != nil {
		return portsout.GitHubManifestConversion{}, normalizeProviderError(err)
	}
	if response.ID < 1 || response.Slug == "" || response.ClientID == "" || response.PEM == "" {
		return portsout.GitHubManifestConversion{}, domain.ErrIntegrationUnavailable
	}
	return portsout.GitHubManifestConversion{AppID: response.ID, AppSlug: response.Slug, AppName: response.Name, ClientID: response.ClientID, PrivateKey: []byte(response.PEM), WebhookSecret: []byte(response.WebhookSecret)}, nil
}
