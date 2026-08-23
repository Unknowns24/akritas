package validationresult

import (
	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
)

func mapError(err error) error {
	return dberrors.ErrValidationResultPersistence.Wrap(err)
}
