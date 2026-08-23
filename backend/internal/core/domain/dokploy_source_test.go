package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewDokployApplicationSource(t *testing.T) {
	serverID := uuid.New()
	source, err := NewDokploySource(serverID, DokploySourceApplication, " app-1 ", "", "api", "API", "production", DokploySourceRunning, "", "")
	if err != nil {
		t.Fatalf("NewDokploySource() error = %v", err)
	}
	if source.ResourceIdentifier != "app-1" || source.ServiceName != "" || source.RuntimeType != "" {
		t.Fatalf("NewDokploySource() = %+v", source)
	}
	if source.IdentityKey() != serverID.String()+":application:app-1:" {
		t.Fatalf("IdentityKey() = %q", source.IdentityKey())
	}
}

func TestNewDokployComposeServiceSource(t *testing.T) {
	serverID := uuid.New()
	source, err := NewDokploySource(serverID, DokploySourceComposeService, " compose-1 ", " api ", "stack-name", "Stack / api", "production", DokploySourceRunning, DokployRuntimeStack, "remote-1")
	if err != nil {
		t.Fatalf("NewDokploySource() error = %v", err)
	}
	if source.ServiceName != "api" || source.RuntimeType != DokployRuntimeStack || source.ProviderServerID != "remote-1" {
		t.Fatalf("NewDokploySource() = %+v", source)
	}
	if source.IdentityKey() != serverID.String()+":compose_service:compose-1:api" {
		t.Fatalf("IdentityKey() = %q", source.IdentityKey())
	}
}

func TestDokploySourceRejectsInvalidDiscriminatorCombinations(t *testing.T) {
	serverID := uuid.New()
	tests := []struct {
		name        string
		typeValue   DokploySourceType
		serviceName string
		runtimeType DokployRuntimeType
	}{
		{name: "application with service", typeValue: DokploySourceApplication, serviceName: "api"},
		{name: "application with runtime", typeValue: DokploySourceApplication, runtimeType: DokployRuntimeCompose},
		{name: "compose without service", typeValue: DokploySourceComposeService, runtimeType: DokployRuntimeCompose},
		{name: "compose without runtime", typeValue: DokploySourceComposeService, serviceName: "api"},
		{name: "unknown type", typeValue: DokploySourceType("database")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewDokploySource(serverID, test.typeValue, "resource", test.serviceName, "instance", "display", "", DokploySourceUnknown, test.runtimeType, "")
			if err == nil {
				t.Fatal("NewDokploySource() expected error")
			}
		})
	}
}

func TestDokploySourceIdentityDistinguishesComposeServices(t *testing.T) {
	serverID := uuid.New()
	api, err := NewDokploySource(serverID, DokploySourceComposeService, "compose", "api", "stack", "API", "", DokploySourceUnknown, DokployRuntimeCompose, "")
	if err != nil {
		t.Fatal(err)
	}
	worker, err := NewDokploySource(serverID, DokploySourceComposeService, "compose", "worker", "stack", "Worker", "", DokploySourceUnknown, DokployRuntimeCompose, "")
	if err != nil {
		t.Fatal(err)
	}
	if api.IdentityKey() == worker.IdentityKey() {
		t.Fatalf("identity collision: %q", api.IdentityKey())
	}
}
