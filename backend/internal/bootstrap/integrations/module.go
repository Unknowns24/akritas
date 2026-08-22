package integrations

import (
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	dokployrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	githubapprepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubapp"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/integrationusage"
	dokployexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/dokploy"
	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/usecase/dokployserver"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/usecase/githubapp"
	"github.com/google/uuid"
)

// Build prepares the integration module only when the administrator middleware
// already exists. cmd/main.go intentionally does not call this before PB-061..063.
func Build(configuration config.Config, admin router.AdminMiddleware) (http.Handler, error) {
	if admin == nil {
		return nil, router.ErrAdminMiddlewareUnavailable
	}
	db, err := postgres.Open(postgres.Config{
		DSN:             configuration.DatabaseURL,
		MaxOpenConns:    configuration.DatabaseMaxOpenConnections,
		MaxIdleConns:    configuration.DatabaseMaxIdleConnections,
		ConnMaxLifetime: configuration.DatabaseConnectionMaxLifetime,
	})
	if err != nil {
		return nil, err
	}
	if err := migrations.Run(db); err != nil {
		return nil, err
	}
	cipher, err := credentialcipher.NewFromKey(configuration.MasterKey)
	if err != nil {
		return nil, err
	}
	credentials, err := credentialstore.New(db, cipher)
	if err != nil {
		return nil, err
	}
	githubAccounts, err := githubrepo.New(db, credentials)
	if err != nil {
		return nil, err
	}
	githubApps, err := githubapprepo.New(db, credentials)
	if err != nil {
		return nil, err
	}
	dokployServers, err := dokployrepo.New(db, credentials)
	if err != nil {
		return nil, err
	}
	githubClient, err := githubexternal.NewClient(githubexternal.ClientConfig{Credentials: credentials, Bindings: githubApps, Now: time.Now})
	if err != nil {
		return nil, err
	}
	dokployClient, err := dokployexternal.NewClient(dokployexternal.ClientConfig{Credentials: credentials})
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
		DokployServers: dokployUseCase, Pagination: pagingConfig, Admin: admin,
	})
}
