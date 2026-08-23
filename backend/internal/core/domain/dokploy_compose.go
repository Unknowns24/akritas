package domain

import (
	"strings"

	"github.com/google/uuid"
)

type DokployCompose struct {
	DokployServerID       uuid.UUID
	ComposeIdentifier     string
	InstanceIdentifier    string
	DisplayName           string
	EnvironmentIdentifier string
	Environment           string
	Status                DokploySourceStatus
	RuntimeType           DokployRuntimeType
	ProviderServerID      string
}

func NewDokployCompose(serverID uuid.UUID, composeID, instanceID, displayName, environmentID, environment string, status DokploySourceStatus, runtimeType DokployRuntimeType, providerServerID string) (DokployCompose, error) {
	value := DokployCompose{
		DokployServerID: serverID, ComposeIdentifier: strings.TrimSpace(composeID),
		InstanceIdentifier: strings.TrimSpace(instanceID), DisplayName: strings.TrimSpace(displayName),
		EnvironmentIdentifier: strings.TrimSpace(environmentID), Environment: strings.TrimSpace(environment),
		Status: status, RuntimeType: runtimeType, ProviderServerID: strings.TrimSpace(providerServerID),
	}
	if err := value.Validate(); err != nil {
		return DokployCompose{}, err
	}
	return value, nil
}

func (c DokployCompose) Validate() error {
	if c.DokployServerID == uuid.Nil || !nonBlank(c.ComposeIdentifier) || !nonBlank(c.InstanceIdentifier) || !nonBlank(c.DisplayName) || c.Status.Validate() != nil {
		return ErrInvalidDokployCompose.Wrap(validationCause("Dokploy compose"))
	}
	if c.RuntimeType != "" && c.RuntimeType.Validate() != nil {
		return ErrInvalidDokployCompose.Wrap(validationCause("Dokploy compose runtime"))
	}
	return nil
}

func (c DokployCompose) Source(serviceName string) (DokploySource, error) {
	serviceName = strings.TrimSpace(serviceName)
	return NewDokploySource(c.DokployServerID, DokploySourceComposeService, c.ComposeIdentifier,
		serviceName, c.InstanceIdentifier, c.DisplayName+" / "+serviceName, c.Environment,
		c.Status, c.RuntimeType, c.ProviderServerID)
}

type DokployComposeService struct {
	DokployServerID   uuid.UUID
	ComposeIdentifier string
	ServiceName       string
}

func NewDokployComposeService(serverID uuid.UUID, composeID, serviceName string) (DokployComposeService, error) {
	value := DokployComposeService{DokployServerID: serverID, ComposeIdentifier: strings.TrimSpace(composeID), ServiceName: strings.TrimSpace(serviceName)}
	if value.DokployServerID == uuid.Nil || !nonBlank(value.ComposeIdentifier) || !nonBlank(value.ServiceName) {
		return DokployComposeService{}, ErrInvalidDokployComposeService.Wrap(validationCause("Dokploy compose service"))
	}
	return value, nil
}
