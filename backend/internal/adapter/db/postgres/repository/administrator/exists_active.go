package administrator

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
)

func (r *repository) ExistsActive(ctx context.Context) (bool, error) {
	var count int64
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Model(&model.Administrator{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
