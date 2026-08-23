package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type DokployServerStore interface {
	CreateWithCredential(context.Context, *domain.DokployServer, SecretValue) error
	Get(context.Context, uuid.UUID) (*domain.DokployServer, error)
	List(context.Context, paging.Params) (paging.Slice[domain.DokployServer], error)
	Update(context.Context, *domain.DokployServer, *SecretValue) error
	UpdateConnection(context.Context, *domain.DokployServer) error
	Delete(context.Context, uuid.UUID) error
}

type DokployValidation struct {
	NormalizedBaseURL string
	ServerIdentifier  string
}

type DokployGateway interface {
	Validate(context.Context, string, string) (DokployValidation, error)
	ValidateUpdate(context.Context, domain.DokployServer, string, *string) (DokployValidation, error)
	TestConnection(context.Context, domain.DokployServer) (ProviderConnectionResult, error)
	ListApplications(context.Context, domain.DokployServer, paging.Params) (paging.Slice[domain.DokployApplication], error)
}
