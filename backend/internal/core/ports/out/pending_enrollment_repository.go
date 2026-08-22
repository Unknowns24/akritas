package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// PendingEnrollmentRepository persists the single pending enrollment slot.
// Save replaces any previous pending enrollment: only one Administrator can
// ever exist, so only one pending enrollment is meaningful at a time.
type PendingEnrollmentRepository interface {
	Save(ctx context.Context, enrollment *domain.PendingEnrollment) error
	// FindByID returns (nil, nil) when no enrollment with that id exists.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.PendingEnrollment, error)
	// Delete consumes the enrollment. Deleting an id that no longer exists
	// is not an error.
	Delete(ctx context.Context, id uuid.UUID) error
}
