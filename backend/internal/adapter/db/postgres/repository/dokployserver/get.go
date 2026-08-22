package dokployserver

import (
	"context"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*domain.DokployServer, error) {
	var server domain.DokployServer
	if err := r.db.WithContext(ctx).Table("dokploy_servers").First(&server, "id = ?", id).Error; err != nil {
		return nil, mapError(err)
	}
	if err := server.Validate(); err != nil {
		return nil, dberrors.ErrIntegrationPersistence.Wrap(err)
	}
	return &server, nil
}
