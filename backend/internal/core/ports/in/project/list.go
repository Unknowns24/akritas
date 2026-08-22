package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

type List interface {
	List(ctx context.Context, query paging.ListQuery) ([]domain.Project, int64, error)
}
