package administratorsession

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func (r *repository) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.AdministratorSession{}).
		Where("id = ?", id).
		UpdateColumn("revoked_at", revokedAt).Error
}
