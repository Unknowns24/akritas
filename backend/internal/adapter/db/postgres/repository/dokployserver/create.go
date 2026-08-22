package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Create(ctx context.Context, server *domain.DokployServer) error {
	return r.db.WithContext(ctx).Create(server).Error
}
