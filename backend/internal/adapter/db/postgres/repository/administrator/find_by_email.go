package administrator

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// FindByEmail returns (nil, nil) when no administrator with that email
// exists.
func (r *repository) FindByEmail(ctx context.Context, email string) (*out.AdministratorAuthentication, error) {
	db := txcontext.From(ctx, r.db).WithContext(ctx)
	var administrator domain.Administrator
	if err := db.Table("administrators").First(&administrator, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := administrator.Validate(); err != nil {
		return nil, err
	}
	var passwordHash string
	if err := db.Table("administrators").Where("id = ?", administrator.ID).Pluck("password_hash", &passwordHash).Error; err != nil {
		return nil, err
	}
	return &out.AdministratorAuthentication{
		Administrator: administrator,
		PasswordHash:  passwordHash,
	}, nil
}
