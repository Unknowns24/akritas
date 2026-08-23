package monitoring

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (s *Service) persistOccurrence(ctx context.Context, project domain.Project, occurrence *domain.PendingLogOccurrence) error {
	exists, err := s.store.OccurrenceExists(ctx, project.ID, occurrence.OccurrenceKey)
	if err != nil || exists {
		return err
	}
	incident, err := s.store.FindGroupableIncident(ctx, project.ID, occurrence.Fingerprint, occurrence.Timestamp, project.MonitoringConfiguration.GroupingWindow)
	if err != nil {
		return err
	}
	if incident == nil {
		key, keyErr := s.store.NextIncidentKey(ctx)
		if keyErr != nil {
			return keyErr
		}
		incident, err = domain.NewIncident(s.newID(), key, project.ID, occurrence.Fingerprint, occurrence.Severity, title(occurrence.Message), occurrence.Timestamp)
		if err != nil {
			return err
		}
		if err = s.store.CreateIncident(ctx, incident); err != nil {
			return err
		}
	} else {
		if err = incident.RecordOccurrence(project.ID, occurrence.Fingerprint, occurrence.Timestamp, project.MonitoringConfiguration.GroupingWindow); err != nil {
			return err
		}
		incident.PromoteSeverity(occurrence.Severity)
		if err = s.store.UpdateIncident(ctx, incident); err != nil {
			return err
		}
	}
	event, err := domain.NewLogEvent(s.newID(), project.ID, occurrence.Timestamp, occurrence.Severity, occurrence.Message, occurrence.Fingerprint, occurrence.DetectionRules, occurrence.ContextBefore, occurrence.ContextAfter)
	if err != nil {
		return err
	}
	if err := event.AssociateOccurrence(incident.ID, project.DokployApplication.ApplicationIdentifier, project.DokployApplication.InstanceIdentifier, occurrence.OccurrenceKey); err != nil {
		return err
	}
	return s.store.CreateLogEvent(ctx, event)
}
