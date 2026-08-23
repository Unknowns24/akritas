package administrator

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *repository) ExistsActive(ctx context.Context) (bool, error) {
	var count int64
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("administrators").Model(&domain.Administrator{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
