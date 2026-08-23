package administrator

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

func (r *repository) RotateCredentials(ctx context.Context, id uuid.UUID, expectedPasswordHash, newPasswordHash string, acceptedTOTPPeriod int64, updatedAt time.Time) (bool, error) {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("administrators").
		Where("id = ? AND password_hash = ?", id, expectedPasswordHash).
		Updates(map[string]any{
			"password_hash":             newPasswordHash,
			"last_accepted_totp_period": acceptedTOTPPeriod,
			"updated_at":                updatedAt,
		})
	return result.RowsAffected == 1, result.Error
}
