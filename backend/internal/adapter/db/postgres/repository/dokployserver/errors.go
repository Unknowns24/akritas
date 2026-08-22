package dokployserver

import (
	"errors"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"gorm.io/gorm"
)

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrIntegrationNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrIntegrationConflict
	}
	return dberrors.ErrIntegrationPersistence.Wrap(err)
}
