package administrator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// FindByID returns (nil, nil) when no administrator with that id exists.
func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Administrator, error) {
	var administrator domain.Administrator
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("administrators").First(&administrator, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := administrator.Validate(); err != nil {
		return nil, err
	}
	return &administrator, nil
}
