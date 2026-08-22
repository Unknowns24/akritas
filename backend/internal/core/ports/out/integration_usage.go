package out

import (
	"context"

	"github.com/google/uuid"
)

type IntegrationUsageReader interface {
	GitHubAccountInUse(context.Context, uuid.UUID) (bool, error)
	DokployServerInUse(context.Context, uuid.UUID) (bool, error)
}
