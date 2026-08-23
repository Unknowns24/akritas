package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type ProjectionStore interface {
	CountGitHubAccounts(context.Context) (int, error)
	CountDokployServers(context.Context) (int, error)
	CountMonitoredProjects(context.Context) (int, error)
	CountActiveIncidents(context.Context) (int, error)
	CountCompletedIncidents(context.Context) (int, error)
	CountPullRequestsCreated(context.Context) (int, error)
	ListActiveIncidents(context.Context, int) ([]domain.Incident, error)
	ListActivity(context.Context, paging.Params) (paging.Slice[domain.ActivityEvent], error)
	FindLastSystemDiagnostics(context.Context) (*domain.Operation, error)
	ListPullRequests(context.Context, paging.Params) (paging.Slice[domain.PullRequestProjection], error)
	GetPullRequest(context.Context, uuid.UUID) (*domain.PullRequestProjection, error)
}
