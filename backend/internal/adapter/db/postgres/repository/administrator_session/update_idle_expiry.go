package administratorsession

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

func (r *repository) UpdateIdleExpiry(ctx context.Context, id uuid.UUID, idleExpiresAt time.Time) error {
	return txcontext.From(ctx, r.db).WithContext(ctx).Table("administrator_sessions").
		Where("id = ?", id).
		UpdateColumn("idle_expires_at", idleExpiresAt).Error
}
