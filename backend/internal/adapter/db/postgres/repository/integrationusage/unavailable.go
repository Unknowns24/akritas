package integrationusage

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// UnavailableReader fails closed until the Project persistence adapter from
// PB-010/PB-011 can answer integration reference checks.
type UnavailableReader struct{}

func (UnavailableReader) GitHubAccountInUse(context.Context, uuid.UUID) (bool, error) {
	return false, domain.ErrIntegrationUnavailable
}

func (UnavailableReader) DokployServerInUse(context.Context, uuid.UUID) (bool, error) {
	return false, domain.ErrIntegrationUnavailable
}
