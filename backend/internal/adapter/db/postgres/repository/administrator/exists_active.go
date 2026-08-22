package administrator

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
)

func (r *repository) ExistsActive(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Administrator{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
