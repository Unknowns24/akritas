package integrations

import (
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	dokployrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	githubapprepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubapp"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/integrationusage"
	dokployexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/dokploy"
	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/usecase/dokployserver"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubapp"
	"github.com/google/uuid"
	"gorm.io/gorm"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type Dependencies struct {
	DB           *gorm.DB
	Credentials  *credentialstore.Store
	Admin        router.AdminMiddleware
	Auth         *authhandler.Handler
	Authenticate portsin.AuthenticateSessionUseCase
}

// Build composes integrations over the application's shared PostgreSQL
// connection and credential store. Infrastructure ownership remains in the
// application bootstrap so auth and integrations cannot silently diverge.
func Build(configuration config.Config, dependencies Dependencies) (http.Handler, error) {
	if dependencies.Admin == nil {
		return nil, router.ErrAdminMiddlewareUnavailable
	}
	if dependencies.DB == nil || dependencies.Credentials == nil || dependencies.Auth == nil || dependencies.Authenticate == nil {
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
	usage := integrationusage.UnavailableReader{}
	githubAccountUseCase := githubaccount.New(githubAccounts, githubClient, usage, uuid.New, time.Now)
	dokployUseCase := dokployserver.New(dokployServers, dokployClient, usage, uuid.New, time.Now)
	githubAppUseCase, err := githubapp.New(githubApps, githubClient, configuration.PublicURL, uuid.New, time.Now)
	if err != nil {
		return nil, err
	}
	pagingConfig, err := pagination.NewConfig(configuration.PaginationSecret, configuration.PaginationTTL)
	if err != nil {
		return nil, err
	}
	return router.New(router.Config{
		GitHubAccounts: githubAccountUseCase, GitHubApps: githubAppUseCase,
		DokployServers: dokployUseCase, Pagination: pagingConfig, Admin: dependencies.Admin,
		Auth: dependencies.Auth, Authenticate: dependencies.Authenticate, AllowedOrigins: configuration.AllowedOrigins,
	})
}
