package administratorsession

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *repository) Save(ctx context.Context, session *domain.AdministratorSession, tokenHash string) error {
	record := map[string]any{
		"id": session.ID, "administrator_id": session.AdministratorID, "token_hash": tokenHash,
		"authenticated_at": session.AuthenticatedAt, "idle_expires_at": session.IdleExpiresAt,
		"absolute_expires_at": session.AbsoluteExpiresAt, "revoked_at": session.RevokedAt,
	}
	return txcontext.From(ctx, r.db).WithContext(ctx).Table("administrator_sessions").Create(record).Error
}
