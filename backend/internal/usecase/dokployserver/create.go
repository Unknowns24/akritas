package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (uc *UseCase) Create(ctx context.Context, command portsin.CreateDokployServerCommand) (*domain.DokployServer, error) {
	validation, err := uc.gateway.Validate(ctx, command.BaseURL, command.APICredential)
	if err != nil {
		return nil, err
	}
	now := uc.now().UTC()
	server, err := domain.NewDokployServer(uc.newID(), command.Name, validation.NormalizedBaseURL, validation.ServerIdentifier, domain.IntegrationStatusConnected, now)
	if err != nil {
		return nil, err
	}
	server.CredentialConfigured = true
	server.LastSyncedAt = &now
	secretBytes := []byte(command.APICredential)
	defer wipe(secretBytes)
	if err := uc.store.CreateWithCredential(ctx, server, portsout.SecretValue{Kind: portsout.SecretKindDokployAPIKey, Plaintext: secretBytes}); err != nil {
		return nil, err
	}
	return server, nil
}
