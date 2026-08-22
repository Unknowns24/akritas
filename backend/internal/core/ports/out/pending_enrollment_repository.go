package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// PendingEnrollmentRepository persists the single pending enrollment slot.
// Save replaces any previous pending enrollment: only one Administrator can
// ever exist, so only one pending enrollment is meaningful at a time.
type PendingEnrollmentRepository interface {
	Save(ctx context.Context, enrollment *domain.PendingEnrollment) error
}
