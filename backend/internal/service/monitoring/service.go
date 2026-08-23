package monitoring

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

const (
	maximumFetchRecords = 10000
	overlapDuration     = time.Second
	pendingTimeout      = 30 * time.Second
)

type Service struct {
	store          portsout.MonitoringStore
	servers        portsout.DokployServerReader
	logs           portsout.LogSource
	transactor     portsout.Transactor
	automation     automationPolicyReader
	investigations investigationStarter
	newID          func() uuid.UUID
	now            func() time.Time
}

type automationPolicyReader interface {
	Get(context.Context) (domain.AutomationPolicy, error)
}

type investigationStarter interface {
	StartIncidentInvestigation(context.Context, portsin.StartIncidentInvestigationCommand) (*domain.Operation, error)
}

func New(store portsout.MonitoringStore, servers portsout.DokployServerReader, logs portsout.LogSource, transactor portsout.Transactor, automation automationPolicyReader, investigations investigationStarter, newID func() uuid.UUID, now func() time.Time) (*Service, error) {
	if store == nil || servers == nil || logs == nil || transactor == nil || newID == nil || now == nil {
		return nil, errors.New("invalid monitoring service configuration")
	}
	return &Service{store: store, servers: servers, logs: logs, transactor: transactor, automation: automation, investigations: investigations, newID: newID, now: now}, nil
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

type automaticInvestigationCandidate struct {
	incidentID  uuid.UUID
	incidentKey string
}

func (s *Service) startAutomaticInvestigations(ctx context.Context, project domain.Project, candidates []automaticInvestigationCandidate) {
	if len(candidates) == 0 || s.automation == nil || s.investigations == nil {
		return
	}
	policy, err := s.automation.Get(ctx)
	if err != nil {
		log.Printf("monitoring: automatic investigation policy lookup failed project_id=%s project_name=%q incidents=%d error=%v", project.ID, project.Name, len(candidates), err)
		return
	}
	if !policy.AutomaticInvestigation {
		log.Printf("monitoring: automatic investigation disabled project_id=%s project_name=%q incidents=%d", project.ID, project.Name, len(candidates))
		return
	}
	for _, candidate := range candidates {
		key := uuid.NewSHA1(uuid.NameSpaceOID, []byte("akritas:automatic-investigation:"+candidate.incidentID.String()))
		operation, startErr := s.investigations.StartIncidentInvestigation(ctx, portsin.StartIncidentInvestigationCommand{
			IncidentID:     candidate.incidentID,
			IdempotencyKey: key,
		})
		if startErr != nil {
			log.Printf("monitoring: automatic investigation failed project_id=%s project_name=%q incident_id=%s incident_key=%s error=%v", project.ID, project.Name, candidate.incidentID, candidate.incidentKey, startErr)
			continue
		}
		if operation == nil {
			log.Printf("monitoring: automatic investigation returned no operation project_id=%s project_name=%q incident_id=%s incident_key=%s", project.ID, project.Name, candidate.incidentID, candidate.incidentKey)
			continue
		}
		log.Printf("monitoring: automatic investigation queued project_id=%s project_name=%q incident_id=%s incident_key=%s operation_id=%s operation_status=%s", project.ID, project.Name, candidate.incidentID, candidate.incidentKey, operation.ID, operation.Status)
	}
}
