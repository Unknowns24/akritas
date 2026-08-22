package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

func (uc *UseCase) List(ctx context.Context, query paging.ListQuery) ([]domain.Project, int64, error) {
	query.Limit = paging.NormalizeLimit(query.Limit)
	if query.Offset < 0 {
		query.Offset = 0
	}
	return uc.projects.List(ctx, query)
}
