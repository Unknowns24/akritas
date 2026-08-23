package dokploy

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) ValidateUpdate(ctx context.Context, server domain.DokployServer, rawURL string, credential *string) (portsout.DokployValidation, error) {
	value := []byte(nil)
	if credential != nil {
		value = []byte(*credential)
	} else if c.credentials != nil {
		stored, err := c.credentials.Get(ctx, portsout.CredentialOwnerDokployServer, server.ID, portsout.SecretKindDokployAPIKey)
		if err != nil {
			return portsout.DokployValidation{}, domain.ErrIntegrationUnavailable
		}
		value = stored
	}
	defer wipe(value)
	return c.Validate(ctx, rawURL, string(value))
}
