package out

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type IncidentStore interface {
	Get(context.Context, uuid.UUID) (*domain.Incident, error)
	List(context.Context, paging.Params) (paging.Slice[domain.Incident], error)
	ListLogEvents(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.LogEvent], error)
}

type MonitoringStore interface {
	ListProjectsForMonitoring(context.Context) ([]domain.Project, error)
	LockProject(context.Context, uuid.UUID) (*domain.Project, error)
	GetCurrentCheckpoint(context.Context, uuid.UUID, bool) (*domain.MonitoringCheckpoint, error)
	CreateCheckpoint(context.Context, *domain.MonitoringCheckpoint) error
	RotateCheckpoint(context.Context, *domain.MonitoringCheckpoint) error
	UpdateCheckpoint(context.Context, *domain.MonitoringCheckpoint, int64) error
	NextIncidentKey(context.Context) (string, error)
	FindGroupableIncident(context.Context, uuid.UUID, string, time.Time, time.Duration) (*domain.Incident, error)
	CreateIncident(context.Context, *domain.Incident) error
	UpdateIncident(context.Context, *domain.Incident) error
	OccurrenceExists(context.Context, uuid.UUID, string) (bool, error)
	CreateLogEvent(context.Context, *domain.LogEvent) error
	UpdateProjectObservation(context.Context, uuid.UUID, domain.MonitoringStatus, time.Time, *time.Time) error
}
