package incident

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// Lock obtains the Incident row for a workflow transition. Callers must invoke
// it inside Transactor.WithinTransaction so the row lock has a bounded lifetime.
func (r *Repository) Lock(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	var incident domain.Incident
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("incidents").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).Take(&incident).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &incident, nil
}
