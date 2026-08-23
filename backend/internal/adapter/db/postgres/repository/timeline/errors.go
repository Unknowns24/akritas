package timeline

import (
	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	return dberrors.ErrIncidentPersistence.Wrap(err)
}

func mapInvalidEvent() error {
	return domain.ErrInvalidIncident
}
