package pendingenrollment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// FindByID returns (nil, nil) when no enrollment with that id exists.
func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*out.PendingEnrollmentAuthentication, error) {
	db := txcontext.From(ctx, r.db).WithContext(ctx)
	var enrollment domain.PendingEnrollment
	if err := db.Table("pending_enrollments").First(&enrollment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := enrollment.Validate(); err != nil {
		return nil, err
	}
	var passwordHash string
	if err := db.Table("pending_enrollments").Where("id = ?", id).Pluck("password_hash", &passwordHash).Error; err != nil {
		return nil, err
	}
	return &out.PendingEnrollmentAuthentication{Enrollment: enrollment, PasswordHash: passwordHash}, nil
}
