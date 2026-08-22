package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// AdministratorSessionRepository persists sessions. tokenHash is the
// persistence-only at-rest representation of the opaque session token; the
// raw token itself is never stored.
type AdministratorSessionRepository interface {
	Save(ctx context.Context, session *domain.AdministratorSession, tokenHash string) error
}
