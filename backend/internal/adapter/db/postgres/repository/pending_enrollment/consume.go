package pendingenrollment

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (r *repository) Consume(ctx context.Context, id uuid.UUID) (*out.PendingEnrollmentAuthentication, error) {
	db := txcontext.From(ctx, r.db).WithContext(ctx)
	var record struct {
		ID           uuid.UUID `gorm:"column:id"`
		Email        string
		DisplayName  string
		PasswordHash string
		CreatedAt    time.Time
		ExpiresAt    time.Time
	}
	result := db.Raw(`DELETE FROM pending_enrollments WHERE id = ? RETURNING id, email, display_name, password_hash, created_at, expires_at`, id).Scan(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 || record.ID == uuid.Nil {
		return nil, nil
	}
	enrollment, err := domain.NewPendingEnrollment(record.ID, record.Email, record.DisplayName, record.CreatedAt, record.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &out.PendingEnrollmentAuthentication{Enrollment: *enrollment, PasswordHash: record.PasswordHash}, nil
}
