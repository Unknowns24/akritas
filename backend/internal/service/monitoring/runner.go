package monitoring

import (
	"context"
	"sync"
	"time"
)

type Runner struct {
	service     *Service
	interval    time.Duration
	concurrency int
}

const maximumRunnerConcurrency = 4

func NewRunner(service *Service, interval time.Duration, concurrency int) (*Runner, error) {
	if service == nil || interval <= 0 || concurrency < 1 || concurrency > maximumRunnerConcurrency {
		return nil, ErrInvalidRunnerConfiguration
	}
	return &Runner{service: service, interval: interval, concurrency: concurrency}, nil
}

func (r *Runner) Run(ctx context.Context) {
	r.runCycle(ctx)
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.runCycle(ctx)
		}
	}
}

func (r *Runner) runCycle(ctx context.Context) {
	projects, err := r.service.store.ListProjectsForMonitoring(ctx)
	if err != nil {
		return
	}
	sem := make(chan struct{}, r.concurrency)
	var group sync.WaitGroup
	for index := range projects {
		if ctx.Err() != nil {
			break
		}
		project := projects[index]
		sem <- struct{}{}
		group.Add(1)
		go func() {
			defer group.Done()
			defer func() { <-sem }()
			_ = r.service.ProcessProject(ctx, project)
		}()
	}
	group.Wait()
}
