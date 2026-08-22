package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type MonitoringStatus string

const (
	MonitoringStatusDisabled   MonitoringStatus = "disabled"
	MonitoringStatusStarting   MonitoringStatus = "starting"
	MonitoringStatusMonitoring MonitoringStatus = "monitoring"
	MonitoringStatusDegraded   MonitoringStatus = "degraded"
	MonitoringStatusError      MonitoringStatus = "error"
)

func (s MonitoringStatus) Validate() error {
	switch s {
	case MonitoringStatusDisabled, MonitoringStatusStarting, MonitoringStatusMonitoring, MonitoringStatusDegraded, MonitoringStatusError:
		return nil
	default:
		return ErrInvalidMonitoringStatus.Wrap(validationCause("monitoring status"))
	}
}

type ProjectHealthStatus string

const (
	ProjectHealthHealthy  ProjectHealthStatus = "healthy"
	ProjectHealthWarning  ProjectHealthStatus = "warning"
	ProjectHealthCritical ProjectHealthStatus = "critical"
	ProjectHealthUnknown  ProjectHealthStatus = "unknown"
)

func (s ProjectHealthStatus) Validate() error {
	switch s {
	case ProjectHealthHealthy, ProjectHealthWarning, ProjectHealthCritical, ProjectHealthUnknown:
		return nil
	default:
		return ErrInvalidProjectHealthStatus.Wrap(validationCause("project health status"))
	}
}

type Project struct {
	ID                      uuid.UUID               `gorm:"type:uuid;primaryKey"`
	Name                    string                  `gorm:"not null"`
	Description             string
	MonitoringStatus        MonitoringStatus        `gorm:"not null"`
	HealthStatus            ProjectHealthStatus     `gorm:"not null"`
	GitHubRepository        GitHubRepository        `gorm:"embedded"`
	DokployApplication      DokployApplication      `gorm:"embedded"`
	MonitoringConfiguration MonitoringConfiguration `gorm:"embedded"`
	LastObservedAt          *time.Time
	CreatedAt               time.Time `gorm:"not null"`
	UpdatedAt               time.Time `gorm:"not null"`
}

func NewProject(
	id uuid.UUID,
	name, description string,
	githubRepository GitHubRepository,
	dokployApplication DokployApplication,
	monitoringConfiguration MonitoringConfiguration,
	createdAt time.Time,
) (*Project, error) {
	monitoringConfiguration.ErrorPatterns = cloneStrings(monitoringConfiguration.ErrorPatterns)
	monitoringConfiguration.IgnoredPatterns = cloneStrings(monitoringConfiguration.IgnoredPatterns)
	monitoringStatus := MonitoringStatusDisabled
	if monitoringConfiguration.Enabled {
		monitoringStatus = MonitoringStatusStarting
	}
	project := &Project{
		ID: id, Name: strings.TrimSpace(name), Description: strings.TrimSpace(description),
		MonitoringStatus: monitoringStatus, HealthStatus: ProjectHealthUnknown,
		GitHubRepository: githubRepository, DokployApplication: dokployApplication,
		MonitoringConfiguration: monitoringConfiguration,
		CreatedAt:               createdAt, UpdatedAt: createdAt,
	}
	if err := project.Validate(); err != nil {
		return nil, err
	}
	if monitoringConfiguration.Enabled {
		if err := project.ReadyForMonitoringEngine(); err != nil {
			return nil, err
		}
	}
	return project, nil
}

func (p Project) Validate() error {
	if p.ID == uuid.Nil || !nonBlank(p.Name) || len(p.Name) > 150 || len(p.Description) > 2000 ||
		p.MonitoringStatus.Validate() != nil || p.HealthStatus.Validate() != nil ||
		p.GitHubRepository.Validate() != nil || p.DokployApplication.Validate() != nil ||
		p.MonitoringConfiguration.Validate() != nil || !validTime(p.CreatedAt) || p.UpdatedAt.Before(p.CreatedAt) {
		return ErrInvalidProject.Wrap(validationCause("project"))
	}
	if p.LastObservedAt != nil && p.LastObservedAt.Before(p.CreatedAt) {
		return ErrInvalidProject.Wrap(validationCause("last observed time"))
	}
	return nil
}

// ReadyForMonitoringEngine reports whether the Project has the non-secret
// GitHub, Dokploy and monitoring snapshots required to start the engine.
// Credentials stay in Credential Store; this does not start workers.
func (p Project) ReadyForMonitoringEngine() error {
	if err := p.GitHubRepository.Validate(); err != nil {
		return err
	}
	if err := p.DokployApplication.Validate(); err != nil {
		return err
	}
	if err := p.MonitoringConfiguration.Validate(); err != nil {
		return err
	}
	if !p.MonitoringConfiguration.Enabled {
		return ErrInvalidMonitoringConfiguration.Wrap(validationCause("monitoring not enabled"))
	}
	return nil
}

func (p *Project) Rename(name, description string, updatedAt time.Time) error {
	if p == nil {
		return ErrInvalidProject.Wrap(validationCause("project"))
	}
	p.Name = strings.TrimSpace(name)
	p.Description = strings.TrimSpace(description)
	p.UpdatedAt = updatedAt
	return p.Validate()
}

func (p *Project) ReplaceIntegrations(repository GitHubRepository, application DokployApplication, updatedAt time.Time) error {
	if p == nil {
		return ErrInvalidProject.Wrap(validationCause("project"))
	}
	p.GitHubRepository = repository
	p.DokployApplication = application
	p.UpdatedAt = updatedAt
	return p.Validate()
}

func (p *Project) ReplaceMonitoringConfiguration(configuration MonitoringConfiguration, updatedAt time.Time) error {
	if p == nil {
		return ErrInvalidProject.Wrap(validationCause("project"))
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	configuration.ErrorPatterns = cloneStrings(configuration.ErrorPatterns)
	configuration.IgnoredPatterns = cloneStrings(configuration.IgnoredPatterns)
	next := *p
	next.MonitoringConfiguration = configuration
	next.UpdatedAt = updatedAt
	if !configuration.Enabled {
		next.MonitoringStatus = MonitoringStatusDisabled
	} else if next.MonitoringStatus == MonitoringStatusDisabled {
		next.MonitoringStatus = MonitoringStatusStarting
	}
	if configuration.Enabled {
		if err := next.ReadyForMonitoringEngine(); err != nil {
			return err
		}
	}
	if err := next.Validate(); err != nil {
		return err
	}
	*p = next
	return nil
}
