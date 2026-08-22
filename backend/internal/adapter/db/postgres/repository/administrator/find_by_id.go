package administrator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// FindByID returns (nil, nil) when no administrator with that id exists.
func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Administrator, error) {
	var record model.Administrator
	if err := txcontext.From(ctx, r.db).WithContext(ctx).First(&record, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return domain.NewAdministrator(record.ID, record.Email, record.DisplayName, record.CreatedAt)
}
