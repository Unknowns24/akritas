package qvac

import (
	"context"
	"errors"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

var ErrInvalidConfiguredRunner = errors.New("invalid configured QVAC runner")

type ClientProvider func(context.Context) (*Client, error)

type ConfiguredRunner struct {
	provider ClientProvider
	tools    *ToolRegistry
	config   RunnerConfig
}

func NewConfiguredRunner(provider ClientProvider, tools *ToolRegistry, config RunnerConfig) (*ConfiguredRunner, error) {
	if provider == nil {
		return nil, ErrInvalidConfiguredRunner
	}
	return &ConfiguredRunner{provider: provider, tools: tools, config: config}, nil
}

func (r *ConfiguredRunner) Run(ctx context.Context, runContext portsout.InvestigationRunContext) (portsout.InvestigationRunResult, error) {
	client, err := r.provider(ctx)
	if err != nil {
		return portsout.InvestigationRunResult{}, err
	}
	config := r.config
	if config.ContextSize <= 0 {
		config.ContextSize = client.ContextSize()
	}
	runner, err := NewRunner(client, r.tools, config)
	if err != nil {
		return portsout.InvestigationRunResult{}, err
	}
	return runner.Run(ctx, runContext)
}
