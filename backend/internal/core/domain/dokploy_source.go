package domain

import (
	"strings"

	"github.com/google/uuid"
)

type DokploySourceType string

const (
	DokploySourceApplication    DokploySourceType = "application"
	DokploySourceComposeService DokploySourceType = "compose_service"
)

func (t DokploySourceType) Validate() error {
	switch t {
	case DokploySourceApplication, DokploySourceComposeService:
		return nil
	default:
		return ErrInvalidDokploySource.Wrap(validationCause("Dokploy source type"))
	}
}

type DokployRuntimeType string

const (
	DokployRuntimeCompose DokployRuntimeType = "docker-compose"
	DokployRuntimeStack   DokployRuntimeType = "stack"
)

func (t DokployRuntimeType) Validate() error {
	switch t {
	case DokployRuntimeCompose, DokployRuntimeStack:
		return nil
	default:
		return ErrInvalidDokploySource.Wrap(validationCause("Dokploy runtime type"))
	}
}

type DokploySourceStatus string

const (
	DokploySourceRunning  DokploySourceStatus = "running"
	DokploySourceStopped  DokploySourceStatus = "stopped"
	DokploySourceDegraded DokploySourceStatus = "degraded"
	DokploySourceUnknown  DokploySourceStatus = "unknown"
)

func (s DokploySourceStatus) Validate() error {
	switch s {
	case DokploySourceRunning, DokploySourceStopped, DokploySourceDegraded, DokploySourceUnknown:
		return nil
	default:
		return ErrInvalidDokploySource.Wrap(validationCause("Dokploy source status"))
	}
}

type DokploySourceSelector struct {
	Type               DokploySourceType
	DokployServerID    uuid.UUID
	ResourceIdentifier string
	ServiceName        string
}

func (s DokploySourceSelector) Normalize() DokploySourceSelector {
	s.ResourceIdentifier = strings.TrimSpace(s.ResourceIdentifier)
	s.ServiceName = strings.TrimSpace(s.ServiceName)
	return s
}

func (s DokploySourceSelector) Validate() error {
	s = s.Normalize()
	if s.Type.Validate() != nil || s.DokployServerID == uuid.Nil || !nonBlank(s.ResourceIdentifier) {
		return ErrInvalidDokploySource.Wrap(validationCause("Dokploy source selector"))
	}
	if (s.Type == DokploySourceApplication && s.ServiceName != "") ||
		(s.Type == DokploySourceComposeService && !nonBlank(s.ServiceName)) {
		return ErrInvalidDokploySource.Wrap(validationCause("Dokploy source selector service"))
	}
	return nil
}

func (s DokploySourceSelector) IdentityKey() string {
	s = s.Normalize()
	return s.DokployServerID.String() + ":" + string(s.Type) + ":" + s.ResourceIdentifier + ":" + s.ServiceName
}

type DokploySource struct {
	Type               DokploySourceType   `gorm:"column:source_type"`
	DokployServerID    uuid.UUID           `gorm:"column:dokploy_server_id;type:uuid"`
	ResourceIdentifier string              `gorm:"column:source_resource_identifier"`
	ServiceName        string              `gorm:"column:source_service_name"`
	InstanceIdentifier string              `gorm:"column:source_instance_identifier"`
	DisplayName        string              `gorm:"column:source_display_name"`
	Environment        string              `gorm:"column:source_environment"`
	Status             DokploySourceStatus `gorm:"column:source_status"`
	RuntimeType        DokployRuntimeType  `gorm:"column:source_runtime_type"`
	ProviderServerID   string              `gorm:"column:source_provider_server_id"`
}

func NewDokploySource(
	dokployServerID uuid.UUID,
	sourceType DokploySourceType,
	resourceIdentifier, serviceName, instanceIdentifier, displayName, environment string,
	status DokploySourceStatus,
	runtimeType DokployRuntimeType,
	providerServerID string,
) (DokploySource, error) {
	source := DokploySource{
		Type: sourceType, DokployServerID: dokployServerID,
		ResourceIdentifier: strings.TrimSpace(resourceIdentifier), ServiceName: strings.TrimSpace(serviceName),
		InstanceIdentifier: strings.TrimSpace(instanceIdentifier), DisplayName: strings.TrimSpace(displayName),
		Environment: strings.TrimSpace(environment), Status: status, RuntimeType: runtimeType,
		ProviderServerID: strings.TrimSpace(providerServerID),
	}
	if err := source.Validate(); err != nil {
		return DokploySource{}, err
	}
	return source, nil
}

func (s DokploySource) Selector() DokploySourceSelector {
	return DokploySourceSelector{Type: s.Type, DokployServerID: s.DokployServerID, ResourceIdentifier: s.ResourceIdentifier, ServiceName: s.ServiceName}
}

func (s DokploySource) IdentityKey() string { return s.Selector().IdentityKey() }

func (s DokploySource) Validate() error {
	if s.Selector().Validate() != nil || !nonBlank(s.InstanceIdentifier) || !nonBlank(s.DisplayName) || s.Status.Validate() != nil {
		return ErrInvalidDokploySource.Wrap(validationCause("Dokploy source"))
	}
	if s.Type == DokploySourceApplication {
		if s.RuntimeType != "" || s.ProviderServerID != "" {
			return ErrInvalidDokploySource.Wrap(validationCause("application source runtime"))
		}
		return nil
	}
	if s.RuntimeType.Validate() != nil {
		return ErrInvalidDokploySource.Wrap(validationCause("compose source runtime"))
	}
	return nil
}

func SourceFromApplication(application DokployApplication) (DokploySource, error) {
	return NewDokploySource(application.DokployServerID, DokploySourceApplication,
		application.ApplicationIdentifier, "", application.InstanceIdentifier, application.DisplayName,
		application.Environment, DokploySourceStatus(application.Status), "", "")
}
