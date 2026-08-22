package administratorsession

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *repository) Save(ctx context.Context, session *domain.AdministratorSession, tokenHash string) error {
	record := &model.AdministratorSession{
		ID:                session.ID,
		AdministratorID:   session.AdministratorID,
		TokenHash:         tokenHash,
		AuthenticatedAt:   session.AuthenticatedAt,
		IdleExpiresAt:     session.IdleExpiresAt,
		AbsoluteExpiresAt: session.AbsoluteExpiresAt,
		RevokedAt:         session.RevokedAt,
	}
	return txcontext.From(ctx, r.db).WithContext(ctx).Create(record).Error
}
