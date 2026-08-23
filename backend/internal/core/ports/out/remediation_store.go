package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// RemediationStore is deliberately minimal (Create+Get only). Full
// lifecycle persistence (status transitions, Changes, PullRequestReference,
// listing) is out of scope for this task and deferred to AKR-55+.
type RemediationStore interface {
	Create(ctx context.Context, value *domain.Remediation) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Remediation, error)
}
