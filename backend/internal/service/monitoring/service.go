package monitoring

import (
	"context"
	"errors"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

const (
	maximumFetchRecords = 10000
	overlapDuration     = time.Second
	pendingTimeout      = 30 * time.Second
)

type Service struct {
	store      portsout.MonitoringStore
	servers    portsout.DokployServerReader
	logs       portsout.LogSource
	transactor portsout.Transactor
	newID      func() uuid.UUID
	now        func() time.Time
}

func New(store portsout.MonitoringStore, servers portsout.DokployServerReader, logs portsout.LogSource, transactor portsout.Transactor, newID func() uuid.UUID, now func() time.Time) (*Service, error) {
	if store == nil || servers == nil || logs == nil || transactor == nil || newID == nil || now == nil {
		return nil, errors.New("invalid monitoring service configuration")
	}
	return &Service{store: store, servers: servers, logs: logs, transactor: transactor, newID: newID, now: now}, nil
}

func (s *Service) PollOnce(ctx context.Context) error {
	projects, err := s.store.ListProjectsForMonitoring(ctx)
	if err != nil {
		return err
	}
	var combined error
	for index := range projects {
		if err := s.ProcessProject(ctx, projects[index]); err != nil {
			combined = errors.Join(combined, err)
		}
	}
	return combined
}
