package administrator

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// FindByEmail returns (nil, nil) when no administrator with that email
// exists.
func (r *repository) FindByEmail(ctx context.Context, email string) (*out.AdministratorCredentials, error) {
	var record model.Administrator
	if err := txcontext.From(ctx, r.db).WithContext(ctx).First(&record, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	admin, err := domain.NewAdministrator(record.ID, record.Email, record.DisplayName, record.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &out.AdministratorCredentials{
		Administrator:          *admin,
		PasswordHash:           record.PasswordHash,
		EncryptedTOTPSecret:    record.EncryptedTOTPSecret,
		LastAcceptedTOTPPeriod: record.LastAcceptedTOTPPeriod,
	}, nil
}
