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
	ID                      uuid.UUID               `gorm:"column:id;type:uuid;primaryKey"`
	Name                    string                  `gorm:"column:name"`
	Description             string                  `gorm:"column:description"`
	MonitoringStatus        MonitoringStatus        `gorm:"column:monitoring_status"`
	HealthStatus            ProjectHealthStatus     `gorm:"column:health_status"`
	GitHubRepository        GitHubRepository        `gorm:"embedded"`
	DokploySource           DokploySource           `gorm:"embedded"`
	MonitoringConfiguration MonitoringConfiguration `gorm:"embedded"`
	LastObservedAt          *time.Time              `gorm:"column:last_observed_at"`
	CreatedAt               time.Time               `gorm:"column:created_at"`
	UpdatedAt               time.Time               `gorm:"column:updated_at"`
}

func NewProject(
	id uuid.UUID,
	name, description string,
	githubRepository GitHubRepository,
	dokploySource DokploySource,
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
		GitHubRepository: githubRepository, DokploySource: dokploySource,
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
		p.GitHubRepository.Validate() != nil || p.DokploySource.Validate() != nil ||
		p.MonitoringConfiguration.Validate() != nil || !validTime(p.CreatedAt) || p.UpdatedAt.Before(p.CreatedAt) {
		return ErrInvalidProject.Wrap(validationCause("project"))
	}
	if p.LastObservedAt != nil && p.LastObservedAt.Before(p.CreatedAt) {
		return ErrInvalidProject.Wrap(validationCause("last observed time"))
	}
	if p.MonitoringConfiguration.Enabled == (p.MonitoringStatus == MonitoringStatusDisabled) {
		return ErrInvalidProject.Wrap(validationCause("monitoring state"))
	}
	return nil
}

func (p Project) ReadyForMonitoringEngine() error {
	if err := p.GitHubRepository.Validate(); err != nil {
		return err
	}
	if err := p.DokploySource.Validate(); err != nil {
		return err
	}
	if err := p.MonitoringConfiguration.Validate(); err != nil {
		return err
	}
	if !p.MonitoringConfiguration.Enabled {
		return ErrInvalidMonitoringConfiguration.Wrap(validationCause("monitoring disabled"))
	}
	return nil
}

func (p *Project) Rename(name, description string, updatedAt time.Time) error {
	if p == nil {
		return ErrInvalidProject
	}
	next := *p
	next.Name = strings.TrimSpace(name)
	next.Description = strings.TrimSpace(description)
	next.UpdatedAt = updatedAt
	if err := next.Validate(); err != nil {
		return err
	}
	*p = next
	return nil
}

func (p *Project) ReplaceIntegrations(repository GitHubRepository, source DokploySource, updatedAt time.Time) error {
	if p == nil {
		return ErrInvalidProject
	}
	if p.MonitoringConfiguration.Enabled {
		return ErrProjectMustBeDisabled
	}
	next := *p
	next.GitHubRepository = repository
	next.DokploySource = source
	next.UpdatedAt = updatedAt
	if err := next.Validate(); err != nil {
		return err
	}
	*p = next
	return nil
}

func (p *Project) RefreshIntegrationSnapshots(repository GitHubRepository, source DokploySource, updatedAt time.Time) error {
	if p == nil || repository.GitHubAccountID != p.GitHubRepository.GitHubAccountID ||
		repository.RepositoryIdentifier != p.GitHubRepository.RepositoryIdentifier ||
		source.IdentityKey() != p.DokploySource.IdentityKey() {
		return ErrInvalidProject.Wrap(validationCause("integration snapshot identity"))
	}
	next := *p
	next.GitHubRepository = repository
	next.DokploySource = source
	next.UpdatedAt = updatedAt
	if err := next.Validate(); err != nil {
		return err
	}
	*p = next
	return nil
}

func (p *Project) ReplaceMonitoringConfiguration(configuration MonitoringConfiguration, updatedAt time.Time) error {
	if p == nil {
		return ErrInvalidProject
	}
	if err := configuration.Validate(); err != nil {
		return err
	}
	next := *p
	next.MonitoringConfiguration = configuration.Clone()
	next.UpdatedAt = updatedAt
	if configuration.Enabled {
		next.MonitoringStatus = MonitoringStatusStarting
		if err := next.ReadyForMonitoringEngine(); err != nil {
			return err
		}
	} else {
		next.MonitoringStatus = MonitoringStatusDisabled
	}
	if err := next.Validate(); err != nil {
		return err
	}
	*p = next
	return nil
}
