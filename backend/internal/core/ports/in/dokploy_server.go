package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type CreateDokployServerCommand struct {
	Name          string
	BaseURL       string
	APICredential string
}

type UpdateDokployServerCommand struct {
	Name          *string
	BaseURL       *string
	APICredential *string
}

type DokployServerUseCase interface {
	Create(context.Context, CreateDokployServerCommand) (*domain.DokployServer, error)
	Get(context.Context, uuid.UUID) (*domain.DokployServer, error)
	List(context.Context, paging.Params) (paging.Slice[domain.DokployServer], error)
	Update(context.Context, uuid.UUID, UpdateDokployServerCommand) (*domain.DokployServer, error)
	Delete(context.Context, uuid.UUID) error
	TestConnection(context.Context, uuid.UUID) (ConnectionTestResult, error)
	ListApplications(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.DokployApplication], error)
}
