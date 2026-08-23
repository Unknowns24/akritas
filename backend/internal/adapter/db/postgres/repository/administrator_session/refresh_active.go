package administratorsession

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *repository) RefreshActive(ctx context.Context, tokenHash string, now, requestedIdleExpiry time.Time) (*domain.AdministratorSession, error) {
	var session domain.AdministratorSession
	result := txcontext.From(ctx, r.db).WithContext(ctx).Raw(`
		UPDATE administrator_sessions
		SET idle_expires_at = LEAST(absolute_expires_at, ?)
		WHERE token_hash = ?
		  AND revoked_at IS NULL
		  AND idle_expires_at > ?
		  AND absolute_expires_at > ?
		RETURNING id, administrator_id, authenticated_at, idle_expires_at, absolute_expires_at, revoked_at`,
		requestedIdleExpiry, tokenHash, now, now,
	).Scan(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if err := session.Validate(); err != nil {
		return nil, err
	}
	return &session, nil
}
