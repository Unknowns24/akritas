package dokployserver

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

func TestCreateValidatesConnectionBeforePersistence(t *testing.T) {
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
	store := &dokployStoreFake{}
	gateway := &dokployGatewayFake{validation: portsout.DokployValidation{NormalizedBaseURL: "https://dokploy.example.com", ServerIdentifier: "fingerprint"}}
	uc := New(store, gateway, usageFake{}, func() uuid.UUID { return uuid.MustParse("a7079eac-151e-49fc-a552-1815dcaa98b9") }, func() time.Time { return now })

	server, err := uc.Create(context.Background(), portsin.CreateDokployServerCommand{Name: " Production ", BaseURL: "https://dokploy.example.com/", APICredential: "dokploy-api-key-value"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !gateway.validated || !store.created || gateway.validationOrder >= store.createOrder {
		t.Fatal("Create() did not validate before persistence")
	}
	if server.Name != "Production" || server.BaseURL != "https://dokploy.example.com" || server.ServerIdentifier != "fingerprint" || server.ConnectionStatus != domain.IntegrationStatusConnected || !server.CredentialConfigured {
		t.Fatalf("Create() server = %+v", server)
	}
}

func TestUpdateRejectedCredentialPreservesOldState(t *testing.T) {
	now := time.Date(2026, 8, 22, 15, 0, 0, 0, time.UTC)
	server, err := domain.NewDokployServer(uuid.New(), "Old", "https://dokploy.example.com", "fingerprint", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server.CredentialConfigured = true
	store := &dokployStoreFake{server: server, secret: "old-key"}
	gateway := &dokployGatewayFake{validationErr: domain.ErrDokployCredentialRejected}
	uc := New(store, gateway, usageFake{}, uuid.New, func() time.Time { return now.Add(time.Minute) })
	name := "New"
	key := "rejected-key-value"

	_, err = uc.Update(context.Background(), server.ID, portsin.UpdateDokployServerCommand{Name: &name, APICredential: &key})
	if !errors.Is(err, domain.ErrDokployCredentialRejected) {
		t.Fatalf("Update() error = %v", err)
	}
	if store.rotated || store.secret != "old-key" || store.server.Name != "Old" {
		t.Fatal("Update() changed persisted state after rejected rotation")
	}
}

type dokployStoreFake struct {
	server      *domain.DokployServer
	secret      string
	created     bool
	createOrder int
	rotated     bool
	deleted     bool
}

func (f *dokployStoreFake) CreateWithCredential(_ context.Context, server *domain.DokployServer, secret portsout.SecretValue) error {
	f.created = true
	f.createOrder = nextOrder()
	copyServer := *server
	f.server = &copyServer
	f.secret = string(secret.Plaintext)
	return nil
}
func (f *dokployStoreFake) Get(_ context.Context, id uuid.UUID) (*domain.DokployServer, error) {
	if f.server == nil || f.server.ID != id {
		return nil, domain.ErrIntegrationNotFound
	}
	copyServer := *f.server
	return &copyServer, nil
}
func (f *dokployStoreFake) List(context.Context, paging.Params) (paging.Slice[domain.DokployServer], error) {
	return paging.Slice[domain.DokployServer]{}, nil
}
func (f *dokployStoreFake) Update(_ context.Context, server *domain.DokployServer, secret *portsout.SecretValue) error {
	copyServer := *server
	f.server = &copyServer
	if secret != nil {
		f.secret = string(secret.Plaintext)
		f.rotated = true
	}
	return nil
}
func (f *dokployStoreFake) Delete(context.Context, uuid.UUID) error { f.deleted = true; return nil }
func (f *dokployStoreFake) UpdateConnection(_ context.Context, server *domain.DokployServer) error {
	copyServer := *server
	f.server = &copyServer
	return nil
}

type dokployGatewayFake struct {
	validation      portsout.DokployValidation
	validationErr   error
	validated       bool
	validationOrder int
}

func (f *dokployGatewayFake) Validate(_ context.Context, _, _ string) (portsout.DokployValidation, error) {
	f.validated = true
	f.validationOrder = nextOrder()
	return f.validation, f.validationErr
}
func (f *dokployGatewayFake) ValidateUpdate(_ context.Context, _ domain.DokployServer, _ string, _ *string) (portsout.DokployValidation, error) {
	f.validated = true
	f.validationOrder = nextOrder()
	return f.validation, f.validationErr
}
func (f *dokployGatewayFake) TestConnection(context.Context, domain.DokployServer) (portsout.ProviderConnectionResult, error) {
	return portsout.ProviderConnectionResult{}, nil
}
func (f *dokployGatewayFake) ListApplications(context.Context, domain.DokployServer, paging.Params) (paging.Slice[domain.DokployApplication], error) {
	return paging.Slice[domain.DokployApplication]{}, nil
}

type usageFake struct{ dokployInUse bool }

func (f usageFake) GitHubAccountInUse(context.Context, uuid.UUID) (bool, error) { return false, nil }
func (f usageFake) DokployServerInUse(context.Context, uuid.UUID) (bool, error) {
	return f.dokployInUse, nil
}

var order int

func nextOrder() int { order++; return order }

var _ portsout.DokployServerStore = (*dokployStoreFake)(nil)
var _ portsout.DokployGateway = (*dokployGatewayFake)(nil)
