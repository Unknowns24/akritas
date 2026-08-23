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
	FindByDokploySource(context.Context, domain.DokploySourceSelector) (*domain.Project, error)
	List(context.Context, paging.Params) (paging.Slice[domain.Project], error)
	Update(context.Context, *domain.Project, time.Time) error
	Delete(context.Context, uuid.UUID) error
}

type GitHubRepositoryResolver interface {
	GetRepository(context.Context, domain.GitHubAccount, string) (domain.GitHubRepository, error)
}

type DokploySourceResolver interface {
	GetApplication(context.Context, domain.DokployServer, string) (domain.DokployApplication, error)
	GetCompose(context.Context, domain.DokployServer, string) (domain.DokployCompose, error)
	ListComposeServices(context.Context, domain.DokployServer, string, bool) ([]domain.DokployComposeService, error)
}

type GitHubAccountReader interface {
	Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error)
}

type DokployServerReader interface {
	Get(context.Context, uuid.UUID) (*domain.DokployServer, error)
}
