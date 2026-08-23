package operation

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"gorm.io/gorm"
)

// FindByIdempotencyKey returns (nil, nil) on a miss: not finding a prior
// operation for this key is the expected, common case, not a failure.
func (r *Repository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Operation, error) {
	var value domain.Operation
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("operations").Where("idempotency_key = ?", key).First(&value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, mapError(err)
	}
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
