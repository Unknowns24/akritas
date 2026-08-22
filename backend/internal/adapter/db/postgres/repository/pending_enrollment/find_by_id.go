package pendingenrollment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// FindByID returns (nil, nil) when no enrollment with that id exists.
func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.PendingEnrollment, error) {
	var record model.PendingEnrollment
	if err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewPendingEnrollment(
		record.ID, record.Email, record.DisplayName, record.PasswordHash,
		record.EncryptedTOTPSecret, record.CreatedAt, record.ExpiresAt,
	)
}
