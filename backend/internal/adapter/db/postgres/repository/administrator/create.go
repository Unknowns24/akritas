package administrator

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// Create maps a unique-email constraint violation to
// domain.ErrAdministratorAlreadyExists instead of propagating the raw
// driver error: two concurrent verifications can both pass the earlier
// ExistsActive check, and only the database constraint catches the race.
func (r *repository) Create(ctx context.Context, administrator *domain.Administrator, passwordHash string) error {
	record := map[string]any{
		"id": administrator.ID, "email": administrator.Email, "display_name": administrator.DisplayName,
		"password_hash": passwordHash, "last_accepted_totp_period": administrator.LastAcceptedTOTPPeriod,
		"created_at": administrator.CreatedAt, "updated_at": administrator.UpdatedAt,
	}
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("administrators").Create(record).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrAdministratorAlreadyExists
		}
		return err
	}
	return nil
}
