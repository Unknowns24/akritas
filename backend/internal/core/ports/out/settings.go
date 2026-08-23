package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type AutomationPolicyStore interface {
	Get(context.Context) (domain.AutomationPolicy, error)
	Put(context.Context, domain.AutomationPolicy) error
}

type QvacConfigurationStore interface {
	Get(context.Context) (domain.QvacConfiguration, error)
	Put(context.Context, domain.QvacConfiguration) error
}
