package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type GitHubIssueReferenceStore interface {
	Create(context.Context, *domain.GitHubIssueReference) error
	FindByInvestigation(context.Context, uuid.UUID) (*domain.GitHubIssueReference, error)
	FindLatestByIncident(context.Context, uuid.UUID) (*domain.GitHubIssueReference, error)
}
