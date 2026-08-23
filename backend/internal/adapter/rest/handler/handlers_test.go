package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type getSetupStatusStub struct{}

func (getSetupStatusStub) Execute(context.Context) (portsin.SetupStatus, error) {
	return portsin.SetupStatus{}, nil
}

type startAdministratorSetupStub struct{}

func (startAdministratorSetupStub) Execute(context.Context, portsin.StartAdministratorSetupInput) (portsin.StartAdministratorSetupOutput, error) {
	return portsin.StartAdministratorSetupOutput{}, nil
}

type verifyAdministratorSetupStub struct{}

func (verifyAdministratorSetupStub) Execute(context.Context, portsin.VerifyAdministratorSetupInput) (portsin.VerifyAdministratorSetupOutput, error) {
	return portsin.VerifyAdministratorSetupOutput{}, nil
}

type loginAdministratorStub struct{}

func (loginAdministratorStub) Execute(context.Context, portsin.LoginAdministratorInput) (portsin.LoginAdministratorOutput, error) {
	return portsin.LoginAdministratorOutput{}, nil
}

type authenticateSessionStub struct{}

func (authenticateSessionStub) Execute(context.Context, string) (domain.AdministratorSession, error) {
	return domain.AdministratorSession{}, nil
}

type getCurrentSessionStub struct{}

func (getCurrentSessionStub) Execute(context.Context, domain.AdministratorSession) (portsin.GetCurrentSessionOutput, error) {
	return portsin.GetCurrentSessionOutput{}, nil
}

type logoutAdministratorStub struct{}

func (logoutAdministratorStub) Execute(context.Context, domain.AdministratorSession) error {
	return nil
}

type githubAccountStub struct{}

func (githubAccountStub) CreatePAT(context.Context, portsin.CreateGitHubPATAccountCommand) (*domain.GitHubAccount, error) {
	return &domain.GitHubAccount{}, nil
}

func (githubAccountStub) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return &domain.GitHubAccount{}, nil
}

func (githubAccountStub) List(context.Context, paging.Params) (paging.Slice[domain.GitHubAccount], error) {
	return paging.Slice[domain.GitHubAccount]{}, nil
}

func (githubAccountStub) Update(context.Context, uuid.UUID, portsin.UpdateGitHubAccountCommand) (*domain.GitHubAccount, error) {
	return &domain.GitHubAccount{}, nil
}

func (githubAccountStub) Delete(context.Context, uuid.UUID) error { return nil }

func (githubAccountStub) TestConnection(context.Context, uuid.UUID) (portsin.ConnectionTestResult, error) {
	return portsin.ConnectionTestResult{}, nil
}

func (githubAccountStub) ListRepositories(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.GitHubRepository], error) {
	return paging.Slice[domain.GitHubRepository]{}, nil
}

type githubAppStub struct{}

func (githubAppStub) StartRegistration(context.Context, portsin.StartGitHubAppRegistrationCommand) (portsin.GitHubManifestRegistrationResult, error) {
	return portsin.GitHubManifestRegistrationResult{}, nil
}

func (githubAppStub) CompleteManifest(context.Context, string, string) (portsin.GitHubManifestCallbackResult, error) {
	return portsin.GitHubManifestCallbackResult{}, nil
}

func (githubAppStub) CompleteInstallation(context.Context, int64, string) (portsin.GitHubInstallationCallbackResult, error) {
	return portsin.GitHubInstallationCallbackResult{}, nil
}

type dokployServerStub struct{}

func (dokployServerStub) Create(context.Context, portsin.CreateDokployServerCommand) (*domain.DokployServer, error) {
	return &domain.DokployServer{}, nil
}

func (dokployServerStub) Get(context.Context, uuid.UUID) (*domain.DokployServer, error) {
	return &domain.DokployServer{}, nil
}

func (dokployServerStub) List(context.Context, paging.Params) (paging.Slice[domain.DokployServer], error) {
	return paging.Slice[domain.DokployServer]{}, nil
}

func (dokployServerStub) Update(context.Context, uuid.UUID, portsin.UpdateDokployServerCommand) (*domain.DokployServer, error) {
	return &domain.DokployServer{}, nil
}

func (dokployServerStub) Delete(context.Context, uuid.UUID) error { return nil }

func (dokployServerStub) TestConnection(context.Context, uuid.UUID) (portsin.ConnectionTestResult, error) {
	return portsin.ConnectionTestResult{}, nil
}

func (dokployServerStub) ListApplications(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.DokployApplication], error) {
	return paging.Slice[domain.DokployApplication]{}, nil
}

type projectStub struct{}

func (projectStub) Create(context.Context, portsin.CreateProjectCommand) (*portsin.ProjectResult, error) {
	return &portsin.ProjectResult{}, nil
}
func (projectStub) Get(context.Context, uuid.UUID) (*portsin.ProjectResult, error) {
	return &portsin.ProjectResult{}, nil
}
func (projectStub) List(context.Context, paging.Params) (paging.Slice[domain.Project], error) {
	return paging.Slice[domain.Project]{}, nil
}
func (projectStub) Update(context.Context, portsin.UpdateProjectCommand) (*portsin.ProjectResult, error) {
	return &portsin.ProjectResult{}, nil
}
func (projectStub) Delete(context.Context, uuid.UUID) error { return nil }
func (projectStub) GetMonitoring(context.Context, uuid.UUID) (domain.MonitoringConfiguration, error) {
	return domain.MonitoringConfiguration{}, nil
}
func (projectStub) PutMonitoring(context.Context, uuid.UUID, domain.MonitoringConfiguration) (domain.MonitoringConfiguration, error) {
	return domain.MonitoringConfiguration{}, nil
}

type incidentStub struct{}

type startAdministratorRecoveryStub struct{}

func (startAdministratorRecoveryStub) Execute(context.Context, portsin.StartAdministratorRecoveryInput) (portsin.StartAdministratorRecoveryOutput, error) {
	return portsin.StartAdministratorRecoveryOutput{}, nil
}

type verifyAdministratorRecoveryStub struct{}

func (verifyAdministratorRecoveryStub) Execute(context.Context, portsin.VerifyAdministratorRecoveryInput) (portsin.VerifyAdministratorRecoveryOutput, error) {
	return portsin.VerifyAdministratorRecoveryOutput{}, nil
}

func (incidentStub) Get(context.Context, uuid.UUID) (*domain.Incident, error) {
	return &domain.Incident{}, nil
}
func (incidentStub) List(context.Context, paging.Params) (paging.Slice[domain.Incident], error) {
	return paging.Slice[domain.Incident]{}, nil
}
func (incidentStub) ListLogEvents(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.LogEvent], error) {
	return paging.Slice[domain.LogEvent]{}, nil
}

func completeUseCases() *portsin.UseCases {
	return &portsin.UseCases{
		GetSetupStatus:              getSetupStatusStub{},
		StartAdministratorSetup:     startAdministratorSetupStub{},
		VerifyAdministratorSetup:    verifyAdministratorSetupStub{},
		LoginAdministrator:          loginAdministratorStub{},
		StartAdministratorRecovery:  startAdministratorRecoveryStub{},
		VerifyAdministratorRecovery: verifyAdministratorRecoveryStub{},
		AuthenticateSession:         authenticateSessionStub{},
		GetCurrentSession:           getCurrentSessionStub{},
		LogoutAdministrator:         logoutAdministratorStub{},
		GitHubAccount:               githubAccountStub{},
		GitHubApp:                   githubAppStub{},
		DokployServer:               dokployServerStub{},
		Project:                     projectStub{},
		Incident:                    incidentStub{},
	}
}

func validHandlersConfig() HandlersConfig {
	return HandlersConfig{
		UseCases: completeUseCases(),
		Pagination: pagination.Config{
			Secret: []byte("01234567890123456789012345678901"),
			TTL:    time.Hour,
		},
		SessionCookieSecure: true,
	}
}

func TestNewHandlersBuildsEveryFeatureHandler(t *testing.T) {
	handlers, err := NewHandlers(validHandlersConfig())
	if err != nil {
		t.Fatalf("NewHandlers() error = %v", err)
	}
	if handlers == nil || handlers.AuthHandler == nil || handlers.GitHubHandler == nil || handlers.DokployHandler == nil || handlers.ProjectHandler == nil || handlers.IncidentHandler == nil {
		t.Fatalf("NewHandlers() = %+v, want every handler", handlers)
	}
}

func TestNewHandlersRejectsIncompleteConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*HandlersConfig)
	}{
		{name: "missing use cases", mutate: func(config *HandlersConfig) { config.UseCases = nil }},
		{name: "missing setup status", mutate: func(config *HandlersConfig) { config.UseCases.GetSetupStatus = nil }},
		{name: "missing start setup", mutate: func(config *HandlersConfig) { config.UseCases.StartAdministratorSetup = nil }},
		{name: "missing verify setup", mutate: func(config *HandlersConfig) { config.UseCases.VerifyAdministratorSetup = nil }},
		{name: "missing login", mutate: func(config *HandlersConfig) { config.UseCases.LoginAdministrator = nil }},
		{name: "missing recovery start", mutate: func(config *HandlersConfig) { config.UseCases.StartAdministratorRecovery = nil }},
		{name: "missing recovery verify", mutate: func(config *HandlersConfig) { config.UseCases.VerifyAdministratorRecovery = nil }},
		{name: "missing authentication", mutate: func(config *HandlersConfig) { config.UseCases.AuthenticateSession = nil }},
		{name: "missing current session", mutate: func(config *HandlersConfig) { config.UseCases.GetCurrentSession = nil }},
		{name: "missing logout", mutate: func(config *HandlersConfig) { config.UseCases.LogoutAdministrator = nil }},
		{name: "missing GitHub account", mutate: func(config *HandlersConfig) { config.UseCases.GitHubAccount = nil }},
		{name: "missing GitHub app", mutate: func(config *HandlersConfig) { config.UseCases.GitHubApp = nil }},
		{name: "missing Dokploy server", mutate: func(config *HandlersConfig) { config.UseCases.DokployServer = nil }},
		{name: "missing Project", mutate: func(config *HandlersConfig) { config.UseCases.Project = nil }},
		{name: "missing Incident", mutate: func(config *HandlersConfig) { config.UseCases.Incident = nil }},
		{name: "invalid pagination", mutate: func(config *HandlersConfig) { config.Pagination = pagination.Config{} }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := validHandlersConfig()
			test.mutate(&config)
			handlers, err := NewHandlers(config)
			if handlers != nil || !errors.Is(err, ErrInvalidHandlersConfiguration) {
				t.Fatalf("NewHandlers() = (%v, %v), want invalid configuration", handlers, err)
			}
		})
	}
}
