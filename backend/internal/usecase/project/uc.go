package project

import (
	"time"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

type UseCase struct {
	projects   portsout.ProjectStore
	accounts   portsout.GitHubAccountReader
	servers    portsout.DokployServerReader
	github     portsout.GitHubRepositoryResolver
	dokploy    portsout.DokploySourceResolver
	monitoring portsout.MonitoringStore
	transactor portsout.Transactor
	newID      func() uuid.UUID
	now        func() time.Time
}

func New(projects portsout.ProjectStore, accounts portsout.GitHubAccountReader, servers portsout.DokployServerReader, github portsout.GitHubRepositoryResolver, dokploy portsout.DokploySourceResolver, newID func() uuid.UUID, now func() time.Time) portsin.ProjectUseCase {
	return &UseCase{projects: projects, accounts: accounts, servers: servers, github: github, dokploy: dokploy, newID: newID, now: now}
}

func NewWithMonitoring(projects portsout.ProjectStore, accounts portsout.GitHubAccountReader, servers portsout.DokployServerReader, github portsout.GitHubRepositoryResolver, dokploy portsout.DokploySourceResolver, monitoring portsout.MonitoringStore, transactor portsout.Transactor, newID func() uuid.UUID, now func() time.Time) portsin.ProjectUseCase {
	return &UseCase{projects: projects, accounts: accounts, servers: servers, github: github, dokploy: dokploy, monitoring: monitoring, transactor: transactor, newID: newID, now: now}
}
