package pendingenrollment

import (
	"context"

	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// Save replaces any previously stored pending enrollment: only one
// Administrator can ever exist, so only one pending enrollment is
// meaningful at a time.
func (r *repository) Save(ctx context.Context, enrollment *domain.PendingEnrollment) error {
	record := toModel(enrollment)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&model.PendingEnrollment{}).Error; err != nil {
			return err
		}
		return tx.Create(record).Error
	})
}

func toModel(enrollment *domain.PendingEnrollment) *model.PendingEnrollment {
	return &model.PendingEnrollment{
		ID:                  enrollment.ID,
		Email:               enrollment.Email,
		DisplayName:         enrollment.DisplayName,
		PasswordHash:        enrollment.PasswordHash,
		EncryptedTOTPSecret: enrollment.EncryptedTOTPSecret,
		CreatedAt:           enrollment.CreatedAt,
		ExpiresAt:           enrollment.ExpiresAt,
	}
}
