package administratorsession

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func (r *repository) UpdateIdleExpiry(ctx context.Context, id uuid.UUID, idleExpiresAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.AdministratorSession{}).
		Where("id = ?", id).
		UpdateColumn("idle_expires_at", idleExpiresAt).Error
}
