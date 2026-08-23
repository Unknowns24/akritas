package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type AutomationUseCase interface {
	GetPolicy(context.Context) (domain.AutomationPolicy, error)
	PutPolicy(context.Context, domain.AutomationPolicy) (domain.AutomationPolicy, error)
}
