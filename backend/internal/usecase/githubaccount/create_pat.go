package githubaccount

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (uc *UseCase) CreatePAT(ctx context.Context, command portsin.CreateGitHubPATAccountCommand) (*domain.GitHubAccount, error) {
	validation, err := uc.gateway.ValidatePAT(ctx, portsout.GitHubPATValidationRequest{
		AccountType: command.AccountType, AccountIdentifier: strings.TrimSpace(command.AccountIdentifier), Token: command.PersonalAccessToken,
	})
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(strings.TrimSpace(validation.AccountIdentifier), strings.TrimSpace(command.AccountIdentifier)) {
		return nil, domain.ErrGitHubCredentialRejected
	}
	if !classicScopesSatisfy(command.AccountType, validation.ClassicScopes) {
		return nil, domain.ErrGitHubCredentialRejected
	}
	now := uc.now().UTC()
	account, err := domain.NewGitHubAccount(uc.newID(), command.DisplayName, command.AccountType, domain.GitHubAuthenticationPersonalAccessToken, command.AccountIdentifier, domain.IntegrationStatusConnected, now)
	if err != nil {
		return nil, err
	}
	account.CredentialConfigured = true
	account.LastCheckedAt = &now
	secretBytes := []byte(command.PersonalAccessToken)
	defer wipe(secretBytes)
	if err := uc.store.CreateWithCredential(ctx, account, portsout.SecretValue{Kind: portsout.SecretKindGitHubPAT, Plaintext: secretBytes}); err != nil {
		return nil, err
	}
	return account, nil
}
