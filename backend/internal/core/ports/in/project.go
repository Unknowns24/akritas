package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type CreateProjectCommand struct {
	Name                    string
	Description             string
	GitHubAccountID         uuid.UUID
	RepositoryIdentifier    string
	DefaultBranch           string
	DokployServerID         uuid.UUID
	ApplicationIdentifier   string
	MonitoringConfiguration domain.MonitoringConfiguration
}

type UpdateProjectCommand struct {
	ID                    uuid.UUID
	Name                  *string
	Description           *string
	GitHubAccountID       *uuid.UUID
	RepositoryIdentifier  *string
	DefaultBranch         *string
	DokployServerID       *uuid.UUID
	ApplicationIdentifier *string
}

type ProjectResult struct {
	Project               *domain.Project
	BuiltInDetectionRules []domain.BuiltInDetectionRule
}

type ProjectUseCase interface {
	Create(context.Context, CreateProjectCommand) (*ProjectResult, error)
	Get(context.Context, uuid.UUID) (*ProjectResult, error)
	List(context.Context, paging.Params) (paging.Slice[domain.Project], error)
	Update(context.Context, UpdateProjectCommand) (*ProjectResult, error)
	Delete(context.Context, uuid.UUID) error
	GetMonitoring(context.Context, uuid.UUID) (domain.MonitoringConfiguration, error)
	PutMonitoring(context.Context, uuid.UUID, domain.MonitoringConfiguration) (domain.MonitoringConfiguration, error)
}
