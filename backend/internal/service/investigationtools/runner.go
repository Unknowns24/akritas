package investigationtools

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/external/qvac"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// Runner resolves the incident repository, attaches read-only GitHub tools, and
// delegates inference to the local QVAC runner.
type Runner struct {
	client   *qvac.Client
	api      GitHubAPI
	resolver *Resolver
	config   qvac.RunnerConfig
}

func NewRunner(client *qvac.Client, api GitHubAPI, resolver *Resolver, config qvac.RunnerConfig) *Runner {
	return &Runner{client: client, api: api, resolver: resolver, config: config}
}

func (r *Runner) Run(ctx context.Context, investigation domain.Investigation) (portsout.InvestigationRunResult, error) {
	scope, err := r.resolver.Resolve(ctx, investigation)
	if err != nil {
		// Without a resolvable repository, still allow structured local inference.
		runner, buildErr := qvac.NewRunner(r.client, qvac.NewToolRegistry(), r.config)
		if buildErr != nil {
			return portsout.InvestigationRunResult{}, buildErr
		}
		return runner.Run(ctx, investigation)
	}
	runner, err := qvac.NewRunner(r.client, Registry(r.api, scope), r.config)
	if err != nil {
		return portsout.InvestigationRunResult{}, err
	}
	return runner.Run(ctx, investigation)
}

var _ portsout.InvestigationRunner = (*Runner)(nil)
