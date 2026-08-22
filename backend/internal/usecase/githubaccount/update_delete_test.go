package githubaccount

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestUpdateRejectedRotationPreservesStoredAccountAndSecret(t *testing.T) {
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
	account, err := domain.NewGitHubAccount(uuid.New(), "Old", domain.GitHubAccountPersonal, domain.GitHubAuthenticationPersonalAccessToken, "unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	account.CredentialConfigured = true
	store := &githubStoreFake{account: account, secret: "old-secret"}
	gateway := &githubGatewayFake{validateErr: domain.ErrGitHubCredentialRejected}
	uc := New(store, gateway, usageFake{}, uuid.New, func() time.Time { return now.Add(time.Minute) })
	name := "New"
	secret := "github_pat_rejected_value_1234567890"

	_, err = uc.Update(context.Background(), account.ID, portsin.UpdateGitHubAccountCommand{DisplayName: &name, PersonalAccessToken: &secret})
	if !errors.Is(err, domain.ErrGitHubCredentialRejected) {
		t.Fatalf("Update() error = %v", err)
	}
	if store.rotated || store.secret != "old-secret" || store.account.DisplayName != "Old" {
		t.Fatal("Update() changed persisted state after rejected rotation")
	}
}

func TestDeleteFailsClosedWhenAccountIsReferenced(t *testing.T) {
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountPersonal, domain.GitHubAuthenticationPersonalAccessToken, "unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	store := &githubStoreFake{account: account}
	uc := New(store, &githubGatewayFake{}, usageFake{githubInUse: true}, uuid.New, time.Now)

	err = uc.Delete(context.Background(), account.ID)
	if !errors.Is(err, domain.ErrIntegrationInUse) {
		t.Fatalf("Delete() error = %v", err)
	}
	if store.deleted {
		t.Fatal("Delete() removed a referenced account")
	}
}

var _ portsout.GitHubAccountStore = (*githubStoreFake)(nil)
var _ portsout.GitHubGateway = (*githubGatewayFake)(nil)
