package monitoring

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (s *Service) recordFailure(ctx context.Context, project domain.Project, cause error) error {
	status := domain.MonitoringStatusDegraded
	if project.MonitoringStatus == domain.MonitoringStatusStarting || project.LastObservedAt == nil {
		status = domain.MonitoringStatusError
	}
	if project.MonitoringConfiguration.Enabled {
		_ = s.transactor.WithinTransaction(ctx, func(txctx context.Context) error {
			if _, err := s.store.LockProject(txctx, project.ID); err != nil {
				return err
			}
			return s.store.UpdateProjectObservation(txctx, project.ID, status, s.now().UTC(), nil)
		})
	}
	return cause
}
