package out

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type ProjectStore interface {
	Create(context.Context, *domain.Project) error
	Get(context.Context, uuid.UUID) (*domain.Project, error)
	FindByNormalizedName(context.Context, string) (*domain.Project, error)
	FindByDokployApplication(context.Context, uuid.UUID, string) (*domain.Project, error)
	List(context.Context, paging.Params) (paging.Slice[domain.Project], error)
	Update(context.Context, *domain.Project, time.Time) error
	Delete(context.Context, uuid.UUID) error
}

type GitHubRepositoryResolver interface {
	GetRepository(context.Context, domain.GitHubAccount, string) (domain.GitHubRepository, error)
}

type DokployApplicationResolver interface {
	GetApplication(context.Context, domain.DokployServer, string) (domain.DokployApplication, error)
}

type GitHubAccountReader interface {
	Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error)
}

type DokployServerReader interface {
	Get(context.Context, uuid.UUID) (*domain.DokployServer, error)
}
