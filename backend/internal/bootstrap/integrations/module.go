package integrations

import (
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	dokployrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	evidencerepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/evidence"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	githubapprepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubapp"
	investigationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/investigation"
	operationrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/operation"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	dokployexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/dokploy"
	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	"github.com/Unknowns24/akritas/backend/internal/adapter/external/qvac"
	resthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/adapter/stub"
	"github.com/Unknowns24/akritas/backend/internal/service/evidenceassembly"
	"github.com/Unknowns24/akritas/backend/internal/service/investigationdispatch"
	"github.com/Unknowns24/akritas/backend/internal/usecase/dokployserver"
	evidenceusecase "github.com/Unknowns24/akritas/backend/internal/usecase/evidence"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubapp"
	investigationusecase "github.com/Unknowns24/akritas/backend/internal/usecase/investigation"
	operationusecase "github.com/Unknowns24/akritas/backend/internal/usecase/operation"
	projectusecase "github.com/Unknowns24/akritas/backend/internal/usecase/project"
	"github.com/google/uuid"
	"gorm.io/gorm"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

// investigationRunTimeout bounds each asynchronous investigation run. It is
// generous for the current stub (which fails immediately) and left as a
// placeholder for a real QVAC call's timeout once PB-028+ lands.
const investigationRunTimeout = 5 * time.Minute

type Dependencies struct {
	DB          *gorm.DB
	Credentials *credentialstore.Store
	Admin       router.AdminMiddleware
	UseCases    *portsin.UseCases
}

// Build composes integrations over the application's shared PostgreSQL
// connection and credential store. Infrastructure ownership remains in the
// application bootstrap so auth and integrations cannot silently diverge.
func Build(configuration config.Config, dependencies Dependencies) (http.Handler, error) {
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
	usage := projects
	githubAccountUseCase := githubaccount.New(githubAccounts, githubClient, usage, uuid.New, time.Now)
	dokployUseCase := dokployserver.New(dokployServers, dokployClient, usage, uuid.New, time.Now)
	projectUseCase := projectusecase.New(projects, githubAccounts, dokployServers, githubClient, dokployClient, uuid.New, time.Now)

	incidentReader := stub.NewDenyAllIncidentReader()
	evidenceAssembler := evidenceassembly.New(incidentReader, projects, uuid.New, time.Now)
	runInvestigation := investigationusecase.NewRunUseCase(investigations, operations, evidences, evidenceAssembler, qvac.NewStubRunner(), time.Now)
	investigationDispatcher := investigationdispatch.New(runInvestigation, investigationRunTimeout)
	investigationUseCase := investigationusecase.New(
		incidentReader, investigations, operations, investigationDispatcher, uuid.New, time.Now,
	)
	operationUseCase := operationusecase.New(operations)
	evidenceUseCase := evidenceusecase.New(investigations, evidences)
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
	useCases.Investigation = investigationUseCase
	useCases.Operation = operationUseCase
	useCases.Evidence = evidenceUseCase
	handlers, err := resthandler.NewHandlers(resthandler.HandlersConfig{
		UseCases:            &useCases,
		Pagination:          pagingConfig,
		SessionCookieSecure: configuration.SessionCookieSecure,
	})
	if err != nil {
		return nil, err
	}
	return router.New(router.Config{
		Handlers: handlers, Admin: dependencies.Admin,
		Authenticate: useCases.AuthenticateSession, AllowedOrigins: configuration.AllowedOrigins,
	})
}
