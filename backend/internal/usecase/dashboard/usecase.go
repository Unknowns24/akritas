package dashboard

import (
	"context"
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

var ErrInvalidUseCase = errors.New("invalid dashboard use case configuration")

type UseCase struct {
	projections portsout.ProjectionStore
}

func New(projections portsout.ProjectionStore) (portsin.DashboardUseCase, error) {
	if projections == nil {
		return nil, ErrInvalidUseCase
	}
	return &UseCase{projections: projections}, nil
}

func (uc *UseCase) GetOverview(ctx context.Context) (domain.Overview, error) {
	monitored, err := uc.projections.CountMonitoredProjects(ctx)
	if err != nil {
		return domain.Overview{}, err
	}
	active, err := uc.projections.CountActiveIncidents(ctx)
	if err != nil {
		return domain.Overview{}, err
	}
	completed, err := uc.projections.CountCompletedIncidents(ctx)
	if err != nil {
		return domain.Overview{}, err
	}
	prs, err := uc.projections.CountPullRequestsCreated(ctx)
	if err != nil {
		return domain.Overview{}, err
	}
	incidents, err := uc.projections.ListActiveIncidents(ctx, 20)
	if err != nil {
		return domain.Overview{}, err
	}
	activity, err := uc.projections.ListActivity(ctx, paging.Params{Limit: 20})
	if err != nil {
		return domain.Overview{}, err
	}
	return domain.Overview{
		MonitoredProjects:          monitored,
		ActiveIncidents:            active,
		WorkflowCompletedIncidents: completed,
		PullRequestsCreated:        prs,
		ActiveInvestigations:       incidents,
		LatestActivity:             activity.Items,
	}, nil
}

func (uc *UseCase) ListActivity(ctx context.Context, params paging.Params) (paging.Slice[domain.ActivityEvent], error) {
	return uc.projections.ListActivity(ctx, params)
}
