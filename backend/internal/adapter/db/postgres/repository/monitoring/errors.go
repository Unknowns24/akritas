package monitoring

import (
	"errors"
	"strings"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrProjectNotFound
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" && strings.Contains(postgresError.ConstraintName, "occurrence") {
		return domain.ErrMonitoringConcurrentModification
	}
	return dberrors.ErrMonitoringPersistence.Wrap(err)
}
