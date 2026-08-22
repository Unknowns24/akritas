package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type PutMonitoring interface {
	PutMonitoring(ctx context.Context, command MonitoringCommand) (domain.MonitoringConfiguration, error)
}
