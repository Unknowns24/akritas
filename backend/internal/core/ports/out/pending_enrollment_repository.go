package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// PendingEnrollmentRepository persists the single pending enrollment slot.
// Save replaces any previous pending enrollment: only one Administrator can
// ever exist, so only one pending enrollment is meaningful at a time.
type PendingEnrollmentAuthentication struct {
	Enrollment   domain.PendingEnrollment
	PasswordHash string
}

type PendingEnrollmentRepository interface {
	// Replace atomically replaces the single enrollment slot and returns the
	// previous owner ID so its credential can be deleted in the same outer transaction.
	Replace(ctx context.Context, enrollment *domain.PendingEnrollment, passwordHash string) (*uuid.UUID, error)
	// FindByID returns (nil, nil) when no enrollment with that id exists.
	FindByID(ctx context.Context, id uuid.UUID) (*PendingEnrollmentAuthentication, error)
	// Delete consumes the enrollment. Deleting an id that no longer exists
	// is not an error.
	Delete(ctx context.Context, id uuid.UUID) error
}
