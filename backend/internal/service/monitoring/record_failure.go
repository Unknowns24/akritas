package monitoring

import (
	"context"
	"log"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (s *Service) recordFailure(ctx context.Context, project domain.Project, cause error) error {
	status := domain.MonitoringStatusDegraded
	if project.MonitoringStatus == domain.MonitoringStatusStarting || project.LastObservedAt == nil {
		status = domain.MonitoringStatusError
	}
	log.Printf("monitoring: project processing failed project_id=%s project_name=%q source_type=%s resource_id=%s service=%s instance=%s previous_status=%s next_status=%s last_observed=%t error=%v", project.ID, project.Name, project.DokploySource.Type, project.DokploySource.ResourceIdentifier, project.DokploySource.ServiceName, project.DokploySource.InstanceIdentifier, project.MonitoringStatus, status, project.LastObservedAt != nil, cause)
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
