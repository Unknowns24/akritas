package monitoring

import (
	"context"
	"log"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (s *Service) persistOccurrence(ctx context.Context, project domain.Project, occurrence *domain.PendingLogOccurrence) (*domain.Incident, bool, error) {
	exists, err := s.store.OccurrenceExists(ctx, project.ID, occurrence.OccurrenceKey)
	if err != nil {
		return nil, false, err
	}
	if exists {
		log.Printf("monitoring: occurrence already persisted project_id=%s project_name=%q occurrence_key=%s", project.ID, project.Name, occurrence.OccurrenceKey)
		return nil, false, nil
	}
	incident, err := s.store.FindGroupableIncident(ctx, project.ID, occurrence.Fingerprint, occurrence.Timestamp, project.MonitoringConfiguration.GroupingWindow)
	if err != nil {
		return nil, false, err
	}
	createdIncident := false
	if incident == nil {
		key, keyErr := s.store.NextIncidentKey(ctx)
		if keyErr != nil {
			return nil, false, keyErr
		}
		incident, err = domain.NewIncident(s.newID(), key, project.ID, occurrence.Fingerprint, occurrence.Severity, title(occurrence.Message), occurrence.Timestamp)
		if err != nil {
			return nil, false, err
		}
		if err = s.store.CreateIncident(ctx, incident); err != nil {
			return nil, false, err
		}
		createdIncident = true
	} else {
		if err = incident.RecordOccurrence(project.ID, occurrence.Fingerprint, occurrence.Timestamp, project.MonitoringConfiguration.GroupingWindow); err != nil {
			return nil, false, err
		}
		incident.PromoteSeverity(occurrence.Severity)
		if err = s.store.UpdateIncident(ctx, incident); err != nil {
			return nil, false, err
		}
	}
	event, err := domain.NewLogEvent(s.newID(), project.ID, occurrence.Timestamp, occurrence.Severity, occurrence.Message, occurrence.Fingerprint, occurrence.DetectionRules, occurrence.ContextBefore, occurrence.ContextAfter)
	if err != nil {
		return nil, false, err
	}
	if err := event.AssociateOccurrence(incident.ID, project.DokploySource, occurrence.OccurrenceKey); err != nil {
		return nil, false, err
	}
	if err := s.store.CreateLogEvent(ctx, event); err != nil {
		return nil, false, err
	}
	log.Printf("monitoring: persisted log occurrence project_id=%s project_name=%q incident_id=%s incident_key=%s created_incident=%t occurrence_key=%s rules=%v severity=%s", project.ID, project.Name, incident.ID, incident.Key, createdIncident, occurrence.OccurrenceKey, occurrence.DetectionRules, occurrence.Severity)
	return incident, createdIncident, nil
}
