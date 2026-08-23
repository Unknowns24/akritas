package monitoring

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (s *Service) ensureCheckpoint(ctx context.Context, project domain.Project) (*domain.MonitoringCheckpoint, error) {
	checkpoint, err := s.store.GetCurrentCheckpoint(ctx, project.ID, false)
	if err != nil || checkpoint != nil {
		return checkpoint, err
	}
	now := s.now().UTC()
	created, err := domain.NewMonitoringCheckpoint(s.newID(), project, domain.InitialLogIngestionFromNow, now)
	if err != nil {
		return nil, err
	}
	err = s.transactor.WithinTransaction(ctx, func(txctx context.Context) error {
		if _, lockErr := s.store.LockProject(txctx, project.ID); lockErr != nil {
			return lockErr
		}
		existing, getErr := s.store.GetCurrentCheckpoint(txctx, project.ID, true)
		if getErr != nil {
			return getErr
		}
		if existing != nil {
			checkpoint = existing
			return nil
		}
		if createErr := s.store.CreateCheckpoint(txctx, created); createErr != nil {
			return createErr
		}
		checkpoint = created
		return nil
	})
	return checkpoint, err
}
