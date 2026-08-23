package administratorsession

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

func (r *repository) RevokeAll(ctx context.Context, administratorID uuid.UUID, revokedAt time.Time) error {
	return txcontext.From(ctx, r.db).WithContext(ctx).Table("administrator_sessions").
		Where("administrator_id = ? AND revoked_at IS NULL", administratorID).
		UpdateColumn("revoked_at", revokedAt).Error
}
