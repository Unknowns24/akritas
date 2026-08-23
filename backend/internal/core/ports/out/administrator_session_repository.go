package out

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// AdministratorSessionRepository persists sessions. tokenHash is the
// persistence-only at-rest representation of the opaque session token; the
// raw token itself is never stored.
type AdministratorSessionRepository interface {
	Save(ctx context.Context, session *domain.AdministratorSession, tokenHash string) error
	// FindByTokenHash returns (nil, nil) when no session with that hash exists.
	FindByTokenHash(ctx context.Context, tokenHash string) (*domain.AdministratorSession, error)
	UpdateIdleExpiry(ctx context.Context, id uuid.UUID, idleExpiresAt time.Time) error
	RefreshActive(ctx context.Context, tokenHash string, now, requestedIdleExpiry time.Time) (*domain.AdministratorSession, error)
	Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
	RevokeAll(ctx context.Context, administratorID uuid.UUID, revokedAt time.Time) error
}
