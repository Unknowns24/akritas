package dokploy

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) TestConnection(ctx context.Context, server domain.DokployServer) (portsout.ProviderConnectionResult, error) {
	startedAt := time.Now()
	result := portsout.ProviderConnectionResult{CheckedAt: startedAt.UTC()}
	if c.credentials == nil {
		result.Status = domain.ConnectionTestUnavailable
		result.UserMessage = "No se pudo acceder a la credencial de Dokploy."
		return result, nil
	}
	credential, err := c.credentials.Get(ctx, portsout.CredentialOwnerDokployServer, server.ID, portsout.SecretKindDokployAPIKey)
	if err != nil {
		result.Status = domain.ConnectionTestUnavailable
		result.UserMessage = "No se pudo acceder a la credencial de Dokploy."
		return result, nil
	}
	defer wipe(credential)
	status, err := c.health(ctx, server.BaseURL, string(credential))
	result.Latency = time.Since(startedAt)
	if err != nil {
		result.Status = domain.ConnectionTestUnavailable
		result.UserMessage = "Dokploy no está disponible temporalmente."
		return result, nil
	}
	result.Status = status
	if status == domain.ConnectionTestConnected {
		result.UserMessage = "La conexión con Dokploy está disponible."
	} else {
		result.UserMessage = "Dokploy rechazó la credencial configurada."
	}
	return result, nil
}
