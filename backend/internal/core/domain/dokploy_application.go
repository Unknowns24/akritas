package domain

import (
	"strings"

	"github.com/google/uuid"
)

type DokployApplicationStatus string

const (
	DokployApplicationRunning  DokployApplicationStatus = "running"
	DokployApplicationStopped  DokployApplicationStatus = "stopped"
	DokployApplicationDegraded DokployApplicationStatus = "degraded"
	DokployApplicationUnknown  DokployApplicationStatus = "unknown"
)

func (s DokployApplicationStatus) Validate() error {
	switch s {
	case DokployApplicationRunning, DokployApplicationStopped, DokployApplicationDegraded, DokployApplicationUnknown:
		return nil
	default:
		return ErrInvalidDokployApplication.Wrap(validationCause("Dokploy application status"))
	}
}

type DokployApplication struct {
	DokployServerID       uuid.UUID
	ApplicationIdentifier string
	InstanceIdentifier    string
	DisplayName           string
	Environment           string
	Status                DokployApplicationStatus
}

func NewDokployApplication(
	dokployServerID uuid.UUID,
	applicationIdentifier, instanceIdentifier, displayName, environment string,
	status DokployApplicationStatus,
) (DokployApplication, error) {
	application := DokployApplication{
		DokployServerID: dokployServerID, ApplicationIdentifier: strings.TrimSpace(applicationIdentifier),
		InstanceIdentifier: strings.TrimSpace(instanceIdentifier), DisplayName: strings.TrimSpace(displayName),
		Environment: strings.TrimSpace(environment), Status: status,
	}
	if err := application.Validate(); err != nil {
		return DokployApplication{}, err
	}
	return application, nil
}

func (a DokployApplication) Validate() error {
	if a.DokployServerID == uuid.Nil || !nonBlank(a.ApplicationIdentifier) || !nonBlank(a.InstanceIdentifier) ||
		!nonBlank(a.DisplayName) || a.Status.Validate() != nil {
		return ErrInvalidDokployApplication.Wrap(validationCause("Dokploy application"))
	}
	return nil
}
