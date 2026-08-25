package qvac

import (
	"context"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestClientUsesStoredContextSize(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	config, err := domain.NewQvacConfigurationWithContext("http://127.0.0.1:11434/v1", 180, 32768, domain.QvacAuthenticationNone, false, "", now)
	if err != nil {
		t.Fatal(err)
	}
	uc, err := New(qvacConfigurationStoreFake{value: config}, qvacCredentialStoreFake{}, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	client, err := uc.Client(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if client.ContextSize() != 32768 {
		t.Fatalf("client context size = %d, want stored value", client.ContextSize())
	}
}

func TestClientUsesStoredLargeContextSize(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	config, err := domain.NewQvacConfigurationWithContext("http://127.0.0.1:11434/v1", 180, 65536, domain.QvacAuthenticationNone, false, "", now)
	if err != nil {
		t.Fatal(err)
	}
	uc, err := New(qvacConfigurationStoreFake{value: config}, qvacCredentialStoreFake{}, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	client, err := uc.Client(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if client.ContextSize() != 65536 {
		t.Fatalf("client context size = %d, want stored value", client.ContextSize())
	}
}

type qvacConfigurationStoreFake struct {
	value domain.QvacConfiguration
}

func (f qvacConfigurationStoreFake) Get(context.Context) (domain.QvacConfiguration, error) {
	return f.value, nil
}

func (f qvacConfigurationStoreFake) Put(context.Context, domain.QvacConfiguration) error {
	return nil
}

type qvacCredentialStoreFake struct{}

func (qvacCredentialStoreFake) Put(context.Context, string, uuid.UUID, portsout.SecretValue) error {
	return nil
}

func (qvacCredentialStoreFake) Get(context.Context, string, uuid.UUID, portsout.SecretKind) ([]byte, error) {
	return nil, nil
}

func (qvacCredentialStoreFake) DeleteOwner(context.Context, string, uuid.UUID) error {
	return nil
}

func (qvacCredentialStoreFake) MoveOwner(context.Context, string, uuid.UUID, string, uuid.UUID) error {
	return nil
}
