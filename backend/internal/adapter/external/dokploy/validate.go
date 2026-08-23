package dokploy

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) Validate(ctx context.Context, rawURL, credential string) (portsout.DokployValidation, error) {
	normalized, err := c.normalizeAndValidateURL(ctx, rawURL)
	if err != nil || strings.TrimSpace(credential) == "" {
		return portsout.DokployValidation{}, domain.ErrDokployCredentialRejected
	}
	status, err := c.health(ctx, normalized, credential)
	if err != nil {
		return portsout.DokployValidation{}, err
	}
	if status != domain.ConnectionTestConnected {
		return portsout.DokployValidation{}, domain.ErrDokployCredentialRejected
	}
	return portsout.DokployValidation{NormalizedBaseURL: normalized, ServerIdentifier: fingerprint(normalized)}, nil
}
