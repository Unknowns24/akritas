package integrations

import (
	"context"
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	dokployrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	evidencerepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/evidence"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	githubapprepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubapp"
	githubissuereferencerepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubissuereference"
	incidentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/incident"
	investigationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/investigation"
	monitoringrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/monitoring"
	operationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/operation"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	timelinerepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/timeline"
	dokployexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/dokploy"
	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	"github.com/Unknowns24/akritas/backend/internal/adapter/external/qvac"
	resthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/service/evidenceassembly"
	"github.com/Unknowns24/akritas/backend/internal/service/investigationdispatch"
	"github.com/Unknowns24/akritas/backend/internal/service/issuecontent"
	monitoringservice "github.com/Unknowns24/akritas/backend/internal/service/monitoring"
	"github.com/Unknowns24/akritas/backend/internal/usecase/dokployserver"
	evidenceusecase "github.com/Unknowns24/akritas/backend/internal/usecase/evidence"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubapp"
	incidentusecase "github.com/Unknowns24/akritas/backend/internal/usecase/incident"
	investigationusecase "github.com/Unknowns24/akritas/backend/internal/usecase/investigation"
	operationusecase "github.com/Unknowns24/akritas/backend/internal/usecase/operation"
	projectusecase "github.com/Unknowns24/akritas/backend/internal/usecase/project"
	"github.com/google/uuid"
	"gorm.io/gorm"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

const investigationRunTimeout = 5 * time.Minute

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

// Build returns only the HTTP handler for callers that do not need direct
// ownership of background runtime components.
func Build(configuration config.Config, dependencies Dependencies) (http.Handler, error) {
	runtime, err := BuildRuntime(configuration, dependencies)
	if err != nil {
		return nil, err
	}

	return runtime.Handler, nil
}

// BuildRuntime composes the integration HTTP runtime, H2 monitoring runner,
// and H3 investigation pipeline over the application's shared infrastructure.
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

	//
	// Persistence
	//

	githubAccounts, err := githubrepo.New(
		dependencies.DB,
		dependencies.Credentials,
	)
	if err != nil {
		return nil, err
	}

	githubApps, err := githubapprepo.New(
		dependencies.DB,
		dependencies.Credentials,
	)
	if err != nil {
		return nil, err
	}

	dokployServers, err := dokployrepo.New(
		dependencies.DB,
		dependencies.Credentials,
	)
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

	investigations, err := investigationrepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}

	operations, err := operationrepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}

	evidences, err := evidencerepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}

	issueReferences, err := githubissuereferencerepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}

	timeline, err := timelinerepo.New(dependencies.DB)
	if err != nil {
		return nil, err
	}

	transactor := postgres.NewTransactor(dependencies.DB)

	//
	// External adapters
	//

	githubClient, err := githubexternal.NewClient(
		githubexternal.ClientConfig{
			Credentials: dependencies.Credentials,
			Bindings:    githubApps,
			Now:         time.Now,
		},
	)
	if err != nil {
		return nil, err
	}

	dokployClient, err := dokployexternal.NewClient(
		dokployexternal.ClientConfig{
			Credentials: dependencies.Credentials,
		},
	)
	if err != nil {
		return nil, err
	}

	qvacClient, err := qvac.NewClient(qvac.ClientConfig{})
	if err != nil {
		return nil, err
	}

	//
	// H1 use cases
	//

	usage := projects

	githubAccountUseCase := githubaccount.New(
		githubAccounts,
		githubClient,
		usage,
		uuid.New,
		time.Now,
	)

	githubAppUseCase, err := githubapp.New(
		githubApps,
		githubClient,
		configuration.PublicURL,
		uuid.New,
		time.Now,
	)
	if err != nil {
		return nil, err
	}

	dokployUseCase := dokployserver.New(
		dokployServers,
		dokployClient,
		usage,
		uuid.New,
		time.Now,
	)

	//
	// H2 — Projects / Monitoring / Incidents
	//

	projectUseCase := projectusecase.NewWithMonitoring(
		projects,
		githubAccounts,
		dokployServers,
		githubClient,
		dokployClient,
		monitoringStore,
		transactor,
		uuid.New,
		time.Now,
	)

	incidentUseCase := incidentusecase.New(incidents, investigations, issueReferences, timeline)

	monitoringService, err := monitoringservice.New(
		monitoringStore,
		dokployServers,
		dokployClient,
		transactor,
		uuid.New,
		time.Now,
	)
	if err != nil {
		return nil, err
	}

	monitoringRunner, err := monitoringservice.NewRunner(
		monitoringService,
		configuration.MonitoringPollInterval,
		configuration.MonitoringConcurrency,
	)
	if err != nil {
		return nil, err
	}

	//
	// H3 — Investigation / Evidence / QVAC
	//
	// Incidents are provided by the real H2 repository.
	//

	evidenceAssembler := evidenceassembly.New(
		incidents,
		projects,
		githubAccounts,
		uuid.New,
		time.Now,
	)

	investigationRunner, err := qvac.NewRunner(
		qvacClient,
		nil,
		qvac.RunnerConfig{RepositoryInspector: githubClient, ContextSize: 16384},
	)
	if err != nil {
		return nil, err
	}

	runInvestigation := investigationusecase.NewRunUseCase(
		incidents,
		investigations,
		operations,
		evidences,
		projects,
		githubAccounts,
		issueReferences,
		githubClient,
		issuecontent.New(),
		evidenceAssembler,
		investigationRunner,
		transactor,
		time.Now,
	)

	investigationDispatcher := investigationdispatch.New(
		runInvestigation,
		investigationRunTimeout,
	)
	if err := runInvestigation.Recover(context.Background(), investigationDispatcher); err != nil {
		return nil, err
	}

	investigationUseCase := investigationusecase.New(
		incidents,
		investigations,
		operations,
		transactor,
		investigationDispatcher,
		uuid.New,
		time.Now,
	)

	operationUseCase := operationusecase.New(operations)

	evidenceUseCase := evidenceusecase.New(
		investigations,
		evidences,
	)

	//
	// REST composition
	//

	pagingConfig, err := pagination.NewConfig(
		configuration.PaginationSecret,
		configuration.PaginationTTL,
	)
	if err != nil {
		return nil, err
	}

	useCases := *dependencies.UseCases

	useCases.GitHubAccount = githubAccountUseCase
	useCases.GitHubApp = githubAppUseCase
	useCases.DokployServer = dokployUseCase
	useCases.Project = projectUseCase
	useCases.Incident = incidentUseCase

	useCases.Investigation = investigationUseCase
	useCases.Operation = operationUseCase
	useCases.Evidence = evidenceUseCase

	handlers, err := resthandler.NewHandlers(
		resthandler.HandlersConfig{
			UseCases:              &useCases,
			Pagination:            pagingConfig,
			SessionCookieSecure:   configuration.SessionCookieSecure,
			SessionCookieSameSite: configuration.SessionCookieSameSite,
		},
	)

	if err != nil {
		return nil, err
	}

	handler, err := router.New(
		router.Config{
			Handlers:       handlers,
			Admin:          dependencies.Admin,
			Authenticate:   useCases.AuthenticateSession,
			AllowedOrigins: configuration.AllowedOrigins,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Handler: handler,
		Monitor: monitoringRunner,
	}, nil
}
