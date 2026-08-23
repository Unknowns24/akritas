// Package investigationdispatch is the minimal InvestigationDispatcher
// implementation: a single goroutine per dispatch, no worker pool or queue.
// It is intentionally minimal until real volume from PB-028+ justifies more.
package investigationdispatch

import (
	"context"
	"time"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type Dispatcher struct {
	run     portsin.RunInvestigationUseCase
	timeout time.Duration
}

func New(run portsin.RunInvestigationUseCase, timeout time.Duration) portsout.InvestigationDispatcher {
	return &Dispatcher{run: run, timeout: timeout}
}

// Dispatch runs on its own context, detached from the request that queued
// the Investigation: the request already returned 202 by the time this runs,
// so it must not inherit that request's cancellation.
func (d *Dispatcher) Dispatch(investigationID, operationID uuid.UUID) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
		defer cancel()
		_ = d.run.Execute(ctx, investigationID, operationID)
	}()
}
