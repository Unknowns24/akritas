package administratorsession

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// FindByTokenHash returns (nil, nil) when no session with that hash exists.
func (r *repository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.AdministratorSession, error) {
	var record model.AdministratorSession
	if err := r.db.WithContext(ctx).First(&record, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	session, err := domain.NewAdministratorSession(
		record.ID, record.AdministratorID, record.AuthenticatedAt, record.IdleExpiresAt, record.AbsoluteExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	session.RevokedAt = record.RevokedAt
	return session, nil
}
