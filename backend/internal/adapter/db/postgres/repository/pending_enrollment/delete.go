package pendingenrollment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

// Delete is idempotent: deleting an id that no longer exists is not an error.
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return txcontext.From(ctx, r.db).WithContext(ctx).Table("pending_enrollments").Where("id = ?", id).Delete(nil).Error
}
