package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	GetByNormalizedName(ctx context.Context, name string) (*domain.Project, error)
	GetByDokployApplication(ctx context.Context, serverID uuid.UUID, applicationIdentifier string) (*domain.Project, error)
	List(ctx context.Context, query paging.ListQuery) ([]domain.Project, int64, error)
	Update(ctx context.Context, project *domain.Project) error
}
