package githubaccount

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, command portsin.UpdateGitHubAccountCommand) (*domain.GitHubAccount, error) {
	current, err := uc.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	candidate := *current
	var secret *portsout.SecretValue
	var secretBytes []byte
	if command.PersonalAccessToken != nil {
		if candidate.AuthenticationMethod != domain.GitHubAuthenticationPersonalAccessToken {
			return nil, domain.ErrIntegrationConflict
		}
		validation, err := uc.gateway.ValidatePAT(ctx, portsout.GitHubPATValidationRequest{
			AccountType: candidate.AccountType, AccountIdentifier: candidate.AccountIdentifier, Token: *command.PersonalAccessToken,
		})
		if err != nil {
			return nil, err
		}
		if !strings.EqualFold(strings.TrimSpace(validation.AccountIdentifier), candidate.AccountIdentifier) || !classicScopesSatisfy(candidate.AccountType, validation.ClassicScopes) {
			return nil, domain.ErrGitHubCredentialRejected
		}
		secretBytes = []byte(*command.PersonalAccessToken)
		defer wipe(secretBytes)
		value := portsout.SecretValue{Kind: portsout.SecretKindGitHubPAT, Plaintext: secretBytes}
		secret = &value
		candidate.AuthenticationStatus = domain.IntegrationStatusConnected
		candidate.CredentialConfigured = true
		checkedAt := uc.now().UTC()
		candidate.LastCheckedAt = &checkedAt
	}
	if command.DisplayName != nil {
		candidate.DisplayName = *command.DisplayName
	}
	candidate.UpdatedAt = uc.now().UTC()
	if err := candidate.Validate(); err != nil {
		return nil, err
	}
	if command.DisplayName == nil && command.PersonalAccessToken == nil {
		return nil, domain.ErrInvalidGitHubAccount
	}
	if err := uc.store.Update(ctx, &candidate, secret); err != nil {
		return nil, err
	}
	return &candidate, nil
}
