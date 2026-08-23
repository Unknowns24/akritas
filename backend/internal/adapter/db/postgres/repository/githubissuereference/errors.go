package githubissuereference

import (
	"errors"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ErrGitHubIssueAlreadyPublished
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return domain.ErrGitHubIssueAlreadyPublished
	}
	return dberrors.ErrGitHubIssueReferencePersistence.Wrap(err)
}
