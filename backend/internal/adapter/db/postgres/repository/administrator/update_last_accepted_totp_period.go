package administrator

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

// UpdateColumn (not Update) deliberately bypasses GORM's auto-updated_at
// convention: domain.NewAdministrator always sets UpdatedAt == CreatedAt,
// so administrators.updated_at must never drift from created_at behind the
// domain's back, or FindByID/FindByEmail would silently reconstruct a
// stale value.
func (r *repository) ConsumeTOTPPeriod(ctx context.Context, id uuid.UUID, period int64) (bool, error) {
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("administrators").
		Where("id = ? AND last_accepted_totp_period < ?", id, period).
		UpdateColumn("last_accepted_totp_period", period)
	return result.RowsAffected == 1, result.Error
}
