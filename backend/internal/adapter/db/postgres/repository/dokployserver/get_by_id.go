package dokployserver

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.DokployServer, error) {
	var server domain.DokployServer
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&server).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrDokployServerNotFound
		}
		return nil, err
	}
	if err := server.Validate(); err != nil {
		return nil, err
	}
	return &server, nil
}
