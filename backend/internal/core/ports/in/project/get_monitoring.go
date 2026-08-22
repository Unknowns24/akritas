package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type GetMonitoring interface {
	GetMonitoring(ctx context.Context, id uuid.UUID) (domain.MonitoringConfiguration, error)
}
