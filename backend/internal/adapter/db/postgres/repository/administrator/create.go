package administrator

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

const uniqueViolationCode = "23505"

// Create maps a unique-email constraint violation to
// domain.ErrAdministratorAlreadyExists instead of propagating the raw
// driver error: two concurrent verifications can both pass the earlier
// ExistsActive check, and only the database constraint catches the race.
func (r *repository) Create(ctx context.Context, administrator *domain.Administrator, passwordHash string, encryptedTOTPSecret []byte) error {
	record := &model.Administrator{
		ID:                  administrator.ID,
		Email:               administrator.Email,
		DisplayName:         administrator.DisplayName,
		PasswordHash:        passwordHash,
		EncryptedTOTPSecret: encryptedTOTPSecret,
		CreatedAt:           administrator.CreatedAt,
		UpdatedAt:           administrator.UpdatedAt,
	}

	if err := txcontext.From(ctx, r.db).WithContext(ctx).Create(record).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return domain.ErrAdministratorAlreadyExists
		}
		return err
	}
	return nil
}
