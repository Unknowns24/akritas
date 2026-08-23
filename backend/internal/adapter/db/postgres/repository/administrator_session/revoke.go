package administratorsession

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

func (r *repository) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	return txcontext.From(ctx, r.db).WithContext(ctx).Table("administrator_sessions").
		Where("id = ? AND revoked_at IS NULL", id).
		UpdateColumn("revoked_at", revokedAt).Error
}
