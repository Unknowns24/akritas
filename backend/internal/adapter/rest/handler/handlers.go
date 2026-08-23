package handler

import (
	"errors"

	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	dokployhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/dokploy"
	evidencehandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/evidence"
	githubhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/github"
	incidenthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/incident"
	investigationhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/investigation"
	operationhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/operation"
	projecthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var ErrInvalidHandlersConfiguration = errors.New("invalid REST handlers configuration")

type Handlers struct {
	AuthHandler          *authhandler.Handler
	GitHubHandler        *githubhandler.Handler
	DokployHandler       *dokployhandler.Handler
	ProjectHandler       *projecthandler.Handler
	IncidentHandler      *incidenthandler.Handler
	InvestigationHandler *investigationhandler.Handler
	OperationHandler     *operationhandler.Handler
	EvidenceHandler      *evidencehandler.Handler
}

type HandlersConfig struct {
	UseCases              *portsin.UseCases
	Pagination            pagination.Config
	SessionCookieSecure   bool
	SessionCookieSameSite string
}

func NewHandlers(config HandlersConfig) (*Handlers, error) {
	if !validUseCases(config.UseCases) {
		return nil, ErrInvalidHandlersConfiguration
	}
	sessionCookieSameSite, err := authhandler.ParseSessionCookieSameSite(config.SessionCookieSameSite)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	githubHandler, err := githubhandler.New(
		config.UseCases.GitHubAccount,
		config.UseCases.GitHubApp,
		config.Pagination,
	)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	dokployHandler, err := dokployhandler.New(config.UseCases.DokployServer, config.Pagination)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	projectHandler, err := projecthandler.New(config.UseCases.Project, config.Pagination)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	incidentHandler, err := incidenthandler.New(config.UseCases.Incident, config.Pagination)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	investigationHandler, err := investigationhandler.New(config.UseCases.Investigation, config.Pagination)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	operationHandler, err := operationhandler.New(config.UseCases.Operation)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	evidenceHandler, err := evidencehandler.New(config.UseCases.Evidence, config.Pagination)
	if err != nil {
		return nil, ErrInvalidHandlersConfiguration
	}

	return &Handlers{
		AuthHandler: authhandler.NewHandlerWithRecovery(
			config.UseCases.GetSetupStatus,
			config.UseCases.StartAdministratorSetup,
			config.UseCases.VerifyAdministratorSetup,
			config.UseCases.LoginAdministrator,
			config.UseCases.StartAdministratorRecovery,
			config.UseCases.VerifyAdministratorRecovery,
			config.UseCases.GetCurrentSession,
			config.UseCases.LogoutAdministrator,
			config.SessionCookieSecure,
			sessionCookieSameSite,
		),
		GitHubHandler:        githubHandler,
		DokployHandler:       dokployHandler,
		ProjectHandler:       projectHandler,
		IncidentHandler:      incidentHandler,
		InvestigationHandler: investigationHandler,
		OperationHandler:     operationHandler,
		EvidenceHandler:      evidenceHandler,
	}, nil
}

func validUseCases(useCases *portsin.UseCases) bool {
	return useCases != nil &&
		useCases.GetSetupStatus != nil &&
		useCases.StartAdministratorSetup != nil &&
		useCases.VerifyAdministratorSetup != nil &&
		useCases.LoginAdministrator != nil &&
		useCases.StartAdministratorRecovery != nil &&
		useCases.VerifyAdministratorRecovery != nil &&
		useCases.AuthenticateSession != nil &&
		useCases.GetCurrentSession != nil &&
		useCases.LogoutAdministrator != nil &&
		useCases.GitHubAccount != nil &&
		useCases.GitHubApp != nil &&
		useCases.DokployServer != nil &&
		useCases.Project != nil &&
		useCases.Incident != nil &&
		useCases.Investigation != nil &&
		useCases.Operation != nil &&
		useCases.Evidence != nil &&
		useCases.Remediation != nil
}
