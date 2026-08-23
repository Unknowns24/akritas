package administratorsession

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// FindByTokenHash returns (nil, nil) when no session with that hash exists.
func (r *repository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.AdministratorSession, error) {
	var session domain.AdministratorSession
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("administrator_sessions").First(&session, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := session.Validate(); err != nil {
		return nil, err
	}
	return &session, nil
}
