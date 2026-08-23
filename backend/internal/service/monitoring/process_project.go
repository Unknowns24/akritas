package monitoring

import (
	"context"
	"reflect"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (s *Service) ProcessProject(ctx context.Context, project domain.Project) error {
	now := s.now().UTC()
	checkpoint, err := s.ensureCheckpoint(ctx, project)
	if err != nil {
		return s.recordFailure(ctx, project, err)
	}
	raw := []portsout.RawLogRecord{}
	if project.MonitoringConfiguration.Enabled {
		server, getErr := s.servers.Get(ctx, project.DokployApplication.DokployServerID)
		if getErr != nil {
			return s.recordFailure(ctx, project, getErr)
		}
		raw, err = s.logs.FetchLogs(ctx, portsout.LogFetchRequest{Server: *server, Application: project.DokployApplication, Tail: maximumFetchRecords, Since: fetchSince(checkpoint, now)})
		if err != nil {
			return s.recordFailure(ctx, project, err)
		}
		raw, err = recordsAfter(checkpoint, raw)
		if err != nil {
			return s.recordFailure(ctx, project, err)
		}
	}
	ready, state := buildDrafts(project, checkpoint, raw, now)
	expectedVersion := checkpoint.Version
	checkpoint.State = state
	checkpoint.NextFinalizeAt = nextFinalizeAt(state.Pending)
	if len(raw) > 0 {
		last := raw[len(raw)-1]
		checkpoint.Advance(domain.MonitoringCursor{Timestamp: last.Timestamp, Ordinal: last.Ordinal, ContentHash: last.ContentHash}, now)
	} else {
		checkpoint.Version++
		checkpoint.UpdatedAt = now
		if checkpoint.InitialBackfillPending && project.MonitoringConfiguration.Enabled {
			checkpoint.InitialBackfillPending = false
		}
	}
	err = s.transactor.WithinTransaction(ctx, func(txctx context.Context) error {
		lockedProject, lockErr := s.store.LockProject(txctx, project.ID)
		if lockErr != nil {
			return lockErr
		}
		lockedCheckpoint, getErr := s.store.GetCurrentCheckpoint(txctx, project.ID, true)
		if getErr != nil {
			return getErr
		}
		if lockedCheckpoint == nil || lockedCheckpoint.ID != checkpoint.ID || lockedCheckpoint.Version != expectedVersion || lockedCheckpoint.SourceApplicationID != project.DokployApplication.ApplicationIdentifier || !reflect.DeepEqual(lockedProject.MonitoringConfiguration, project.MonitoringConfiguration) {
			return domain.ErrMonitoringConcurrentModification
		}
		for index := range ready {
			if persistErr := s.persistOccurrence(txctx, *lockedProject, &ready[index]); persistErr != nil {
				return persistErr
			}
		}
		if updateErr := s.store.UpdateCheckpoint(txctx, checkpoint, expectedVersion); updateErr != nil {
			return updateErr
		}
		status := domain.MonitoringStatusDisabled
		var observedAt *time.Time
		if lockedProject.MonitoringConfiguration.Enabled {
			status = domain.MonitoringStatusMonitoring
			observedAt = &now
		}
		return s.store.UpdateProjectObservation(txctx, project.ID, status, now, observedAt)
	})
	if err != nil {
		return s.recordFailure(ctx, project, err)
	}
	return nil
}

func nextFinalizeAt(values []domain.PendingLogOccurrence) *time.Time {
	if len(values) == 0 {
		return nil
	}
	next := values[0].FinalizeAfter
	for _, value := range values[1:] {
		if value.FinalizeAfter.Before(next) {
			next = value.FinalizeAfter
		}
	}
	return &next
}
