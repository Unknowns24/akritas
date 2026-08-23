package integrations

import (
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	dokployrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	githubapprepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubapp"
	incidentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/incident"
	monitoringrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/monitoring"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	dokployexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/dokploy"
	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	resthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	monitoringservice "github.com/Unknowns24/akritas/backend/internal/service/monitoring"
	"github.com/Unknowns24/akritas/backend/internal/usecase/dokployserver"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubapp"
	incidentusecase "github.com/Unknowns24/akritas/backend/internal/usecase/incident"
	projectusecase "github.com/Unknowns24/akritas/backend/internal/usecase/project"
	"github.com/google/uuid"
	"gorm.io/gorm"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type Dependencies struct {
	DB          *gorm.DB
	Credentials *credentialstore.Store
	Admin       router.AdminMiddleware
	UseCases    *portsin.UseCases
}

type Runtime struct {
	Handler http.Handler
	Monitor *monitoringservice.Runner
}

// Build composes integrations over the application's shared PostgreSQL
// connection and credential store. Infrastructure ownership remains in the
// application bootstrap so auth and integrations cannot silently diverge.
func Build(configuration config.Config, dependencies Dependencies) (http.Handler, error) {
	runtime, err := BuildRuntime(configuration, dependencies)
	if err != nil {
		return nil, err
	}
	return runtime.Handler, nil
}

func BuildRuntime(configuration config.Config, dependencies Dependencies) (*Runtime, error) {
	if dependencies.Admin == nil {
		return nil, router.ErrAdminMiddlewareUnavailable
	}
	if dependencies.UseCases == nil {
		return nil, resthandler.ErrInvalidHandlersConfiguration
	}
	if dependencies.DB == nil || dependencies.Credentials == nil {
		return nil, router.ErrInvalidRouterConfiguration
	}
	githubAccounts, err := githubrepo.New(dependencies.DB, dependencies.Credentials)
	if err != nil {
		return nil, err
	}
	githubApps, err := githubapprepo.New(dependencies.DB, dependencies.Credentials)
	if err != nil {
		return nil, err
	}
	dokployServers, err := dokployrepo.New(dependencies.DB, dependencies.Credentials)
	if err != nil {
		return nil, err
	}
	githubClient, err := githubexternal.NewClient(githubexternal.ClientConfig{Credentials: dependencies.Credentials, Bindings: githubApps, Now: time.Now})
	if err != nil {
		return nil, err
	}
	dokployClient, err := dokployexternal.NewClient(dokployexternal.ClientConfig{Credentials: dependencies.Credentials})
	if err != nil {
		return nil, err
	}
	projects, err := projectrepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}
	incidents, err := incidentrepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}
	monitoringStore, err := monitoringrepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}
	transactor := postgres.NewTransactor(dependencies.DB)
	usage := projects
	githubAccountUseCase := githubaccount.New(githubAccounts, githubClient, usage, uuid.New, time.Now)
	dokployUseCase := dokployserver.New(dokployServers, dokployClient, usage, uuid.New, time.Now)
	projectUseCase := projectusecase.NewWithMonitoring(projects, githubAccounts, dokployServers, githubClient, dokployClient, monitoringStore, transactor, uuid.New, time.Now)
	incidentUseCase := incidentusecase.New(incidents)
	monitor, err := monitoringservice.New(monitoringStore, dokployServers, dokployClient, transactor, uuid.New, time.Now)
	if err != nil {
		return nil, err
	}
	runner, err := monitoringservice.NewRunner(monitor, configuration.MonitoringPollInterval, configuration.MonitoringConcurrency)
	if err != nil {
		return nil, err
	}
	githubAppUseCase, err := githubapp.New(githubApps, githubClient, configuration.PublicURL, uuid.New, time.Now)
	if err != nil {
		return nil, err
	}
	pagingConfig, err := pagination.NewConfig(configuration.PaginationSecret, configuration.PaginationTTL)
	if err != nil {
		return nil, err
	}
	useCases := *dependencies.UseCases
	useCases.GitHubAccount = githubAccountUseCase
	useCases.GitHubApp = githubAppUseCase
	useCases.DokployServer = dokployUseCase
	useCases.Project = projectUseCase
	useCases.Incident = incidentUseCase
	handlers, err := resthandler.NewHandlers(resthandler.HandlersConfig{
		UseCases:            &useCases,
		Pagination:          pagingConfig,
		SessionCookieSecure: configuration.SessionCookieSecure,
	})
	if err != nil {
		return nil, err
	}
	handler, err := router.New(router.Config{
		Handlers: handlers, Admin: dependencies.Admin,
		Authenticate: useCases.AuthenticateSession, AllowedOrigins: configuration.AllowedOrigins,
	})
	if err != nil {
		return nil, err
	}
	return &Runtime{Handler: handler, Monitor: runner}, nil
}
