package pendingenrollment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

// Delete is idempotent: deleting an id that no longer exists is not an error.
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.PendingEnrollment{}, "id = ?", id).Error
}
