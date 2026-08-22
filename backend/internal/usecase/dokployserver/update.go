package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, command portsin.UpdateDokployServerCommand) (*domain.DokployServer, error) {
	current, err := uc.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	candidate := *current
	baseURL := current.BaseURL
	if command.BaseURL != nil {
		baseURL = *command.BaseURL
	}
	var secret *portsout.SecretValue
	var secretBytes []byte
	if command.BaseURL != nil || command.APICredential != nil {
		validation, err := uc.gateway.ValidateUpdate(ctx, *current, baseURL, command.APICredential)
		if err != nil {
			return nil, err
		}
		candidate.BaseURL = validation.NormalizedBaseURL
		candidate.ServerIdentifier = validation.ServerIdentifier
		candidate.ConnectionStatus = domain.IntegrationStatusConnected
		candidate.CredentialConfigured = true
		checkedAt := uc.now().UTC()
		candidate.LastSyncedAt = &checkedAt
		if command.APICredential != nil {
			secretBytes = []byte(*command.APICredential)
			defer wipe(secretBytes)
			value := portsout.SecretValue{Kind: portsout.SecretKindDokployAPIKey, Plaintext: secretBytes}
			secret = &value
		}
	}
	if command.Name != nil {
		candidate.Name = *command.Name
	}
	if command.Name == nil && command.BaseURL == nil && command.APICredential == nil {
		return nil, domain.ErrInvalidDokployServer
	}
	candidate.UpdatedAt = uc.now().UTC()
	if err := candidate.Validate(); err != nil {
		return nil, err
	}
	if err := uc.store.Update(ctx, &candidate, secret); err != nil {
		return nil, err
	}
	return &candidate, nil
}
