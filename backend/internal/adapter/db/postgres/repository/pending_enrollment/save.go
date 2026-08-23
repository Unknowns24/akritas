package pendingenrollment

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// Replace replaces any previously stored pending enrollment: only one
// Administrator can ever exist, so only one pending enrollment is
// meaningful at a time.
func (r *repository) Replace(ctx context.Context, enrollment *domain.PendingEnrollment, passwordHash string) (*uuid.UUID, error) {
	var previousID *uuid.UUID
	db := txcontext.From(ctx, r.db).WithContext(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE pending_enrollments IN EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		var ids []uuid.UUID
		if err := tx.Table("pending_enrollments").Pluck("id", &ids).Error; err != nil {
			return err
		}
		if len(ids) > 0 {
			id := ids[0]
			previousID = &id
		}
		if err := tx.Table("pending_enrollments").Where("1 = 1").Delete(nil).Error; err != nil {
			return err
		}
		return tx.Table("pending_enrollments").Create(map[string]any{
			"id": enrollment.ID, "email": enrollment.Email, "display_name": enrollment.DisplayName,
			"password_hash": passwordHash, "created_at": enrollment.CreatedAt, "expires_at": enrollment.ExpiresAt,
		}).Error
	})
	return previousID, err
}
