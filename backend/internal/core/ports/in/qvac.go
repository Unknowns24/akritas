package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type QvacAuthenticationCommand struct {
	Type          domain.QvacAuthenticationType
	BearerToken   string
	BasicUsername string
	BasicPassword string
}

type PutQvacConfigurationCommand struct {
	EndpointURL              string
	ConnectionTimeoutSeconds int
	ContextSize              int
	Authentication           QvacAuthenticationCommand
}

type QvacRuntimeStatus struct {
	ConnectionStatus domain.IntegrationStatus
	Runtime          string
	ActiveModel      string
	Version          string
	Latency          int64
	CheckedAt        string
}

type QvacUseCase interface {
	GetConfiguration(context.Context) (domain.QvacConfiguration, error)
	PutConfiguration(context.Context, PutQvacConfigurationCommand) (domain.QvacConfiguration, error)
	TestConnection(context.Context) (ConnectionTestResult, error)
	GetStatus(context.Context) (QvacRuntimeStatus, error)
}
