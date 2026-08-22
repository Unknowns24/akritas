package githubaccount

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func TestCreatePATValidatesBeforePersistenceAndReturnsSafeProjection(t *testing.T) {
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
	store := &githubStoreFake{}
	gateway := &githubGatewayFake{validateResult: portsout.GitHubPATValidation{AccountIdentifier: "Unknowns24"}}
	uc := New(store, gateway, usageFake{}, func() uuid.UUID { return uuid.MustParse("b8d70b4e-e96d-45f7-a50f-f8f896385c43") }, func() time.Time { return now })

	account, err := uc.CreatePAT(context.Background(), portsin.CreateGitHubPATAccountCommand{
		DisplayName:         " Akritas ",
		AccountType:         domain.GitHubAccountPersonal,
		AccountIdentifier:   "unknowns24",
		PersonalAccessToken: "github_pat_secret_value_1234567890",
	})
	if err != nil {
		t.Fatalf("CreatePAT() error = %v", err)
	}
	if !gateway.validated || !store.created || gateway.validationOrder >= store.createOrder {
		t.Fatal("CreatePAT() did not validate before persistence")
	}
	if store.secret != "github_pat_secret_value_1234567890" {
		t.Fatal("CreatePAT() did not pass the write-only secret to secure persistence")
	}
	if account.DisplayName != "Akritas" || account.AuthenticationMethod != domain.GitHubAuthenticationPersonalAccessToken || account.AuthenticationStatus != domain.IntegrationStatusConnected || !account.CredentialConfigured {
		t.Fatalf("CreatePAT() account = %+v", account)
	}
}

func TestCreatePATDoesNotPersistRejectedCredential(t *testing.T) {
	store := &githubStoreFake{}
	gateway := &githubGatewayFake{validateErr: domain.ErrGitHubCredentialRejected}
	uc := New(store, gateway, usageFake{}, uuid.New, time.Now)

	_, err := uc.CreatePAT(context.Background(), portsin.CreateGitHubPATAccountCommand{
		DisplayName: "Akritas", AccountType: domain.GitHubAccountPersonal,
		AccountIdentifier: "unknowns24", PersonalAccessToken: "github_pat_invalid_value_1234567890",
	})
	if !errors.Is(err, domain.ErrGitHubCredentialRejected) {
		t.Fatalf("CreatePAT() error = %v", err)
	}
	if store.created {
		t.Fatal("CreatePAT() persisted a rejected credential")
	}
}

type githubStoreFake struct {
	account     *domain.GitHubAccount
	secret      string
	created     bool
	createOrder int
	deleted     bool
	rotated     bool
}

func (f *githubStoreFake) CreateWithCredential(_ context.Context, account *domain.GitHubAccount, secret portsout.SecretValue) error {
	f.created = true
	f.createOrder = nextOrder()
	copyAccount := *account
	f.account = &copyAccount
	f.secret = string(secret.Plaintext)
	return nil
}

func (f *githubStoreFake) Get(_ context.Context, id uuid.UUID) (*domain.GitHubAccount, error) {
	if f.account == nil || f.account.ID != id {
		return nil, domain.ErrIntegrationNotFound
	}
	copyAccount := *f.account
	return &copyAccount, nil
}

func (f *githubStoreFake) List(_ context.Context, _ paging.Params) (paging.Slice[domain.GitHubAccount], error) {
	if f.account == nil {
		return paging.Slice[domain.GitHubAccount]{}, nil
	}
	return paging.Slice[domain.GitHubAccount]{Items: []domain.GitHubAccount{*f.account}, Total: 1}, nil
}

func (f *githubStoreFake) Update(_ context.Context, account *domain.GitHubAccount, secret *portsout.SecretValue) error {
	copyAccount := *account
	f.account = &copyAccount
	if secret != nil {
		f.secret = string(secret.Plaintext)
		f.rotated = true
	}
	return nil
}

func (f *githubStoreFake) Delete(_ context.Context, _ uuid.UUID) error {
	f.deleted = true
	return nil
}

func (f *githubStoreFake) UpdateConnection(_ context.Context, account *domain.GitHubAccount) error {
	copyAccount := *account
	f.account = &copyAccount
	return nil
}

type githubGatewayFake struct {
	validateResult  portsout.GitHubPATValidation
	validateErr     error
	validated       bool
	validationOrder int
	connection      portsout.ProviderConnectionResult
	repositories    paging.Slice[domain.GitHubRepository]
}

func (f *githubGatewayFake) ValidatePAT(_ context.Context, _ portsout.GitHubPATValidationRequest) (portsout.GitHubPATValidation, error) {
	f.validated = true
	f.validationOrder = nextOrder()
	return f.validateResult, f.validateErr
}

func (f *githubGatewayFake) TestConnection(_ context.Context, _ domain.GitHubAccount) (portsout.ProviderConnectionResult, error) {
	return f.connection, nil
}

func (f *githubGatewayFake) ListRepositories(_ context.Context, _ domain.GitHubAccount, _ paging.Params) (paging.Slice[domain.GitHubRepository], error) {
	return f.repositories, nil
}

var operationOrder int

func nextOrder() int {
	operationOrder++
	return operationOrder
}

type usageFake struct {
	githubInUse  bool
	dokployInUse bool
}

func (f usageFake) GitHubAccountInUse(context.Context, uuid.UUID) (bool, error) {
	return f.githubInUse, nil
}
func (f usageFake) DokployServerInUse(context.Context, uuid.UUID) (bool, error) {
	return f.dokployInUse, nil
}
