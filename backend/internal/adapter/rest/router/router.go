package router

import (
	"errors"
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/github"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var (
	ErrAdminMiddlewareUnavailable = errors.New("administrator middleware is unavailable")
	ErrInvalidRouterConfiguration = errors.New("invalid integration router configuration")
)

type AdminMiddleware func(http.Handler) http.Handler

type Config struct {
	GitHubAccounts portsin.GitHubAccountUseCase
	GitHubApps     portsin.GitHubAppUseCase
	DokployServers portsin.DokployServerUseCase
	Pagination     pagination.Config
	Admin          AdminMiddleware
}

func New(config Config) (http.Handler, error) {
	if config.Admin == nil {
		return nil, ErrAdminMiddlewareUnavailable
	}
	githubHandler, err := github.New(config.GitHubAccounts, config.GitHubApps, config.Pagination)
	if err != nil {
		return nil, ErrInvalidRouterConfiguration
	}
	dokployHandler, err := dokploy.New(config.DokployServers, config.Pagination)
	if err != nil {
		return nil, ErrInvalidRouterConfiguration
	}
	private := http.NewServeMux()
	private.HandleFunc("GET /api/v1/integrations/github/accounts", githubHandler.ListAccounts)
	private.HandleFunc("POST /api/v1/integrations/github/accounts", githubHandler.CreateAccount)
	private.HandleFunc("GET /api/v1/integrations/github/accounts/{account_id}", githubHandler.GetAccount)
	private.HandleFunc("PATCH /api/v1/integrations/github/accounts/{account_id}", githubHandler.UpdateAccount)
	private.HandleFunc("DELETE /api/v1/integrations/github/accounts/{account_id}", githubHandler.DeleteAccount)
	private.HandleFunc("POST /api/v1/integrations/github/accounts/{account_id}/connection-test", githubHandler.TestConnection)
	private.HandleFunc("GET /api/v1/integrations/github/accounts/{account_id}/repositories", githubHandler.ListRepositories)
	private.HandleFunc("POST /api/v1/integrations/github/app-manifest/registrations", githubHandler.StartManifest)
	private.HandleFunc("GET /api/v1/integrations/dokploy/servers", dokployHandler.List)
	private.HandleFunc("POST /api/v1/integrations/dokploy/servers", dokployHandler.Create)
	private.HandleFunc("GET /api/v1/integrations/dokploy/servers/{server_id}", dokployHandler.Get)
	private.HandleFunc("PATCH /api/v1/integrations/dokploy/servers/{server_id}", dokployHandler.Update)
	private.HandleFunc("DELETE /api/v1/integrations/dokploy/servers/{server_id}", dokployHandler.Delete)
	private.HandleFunc("POST /api/v1/integrations/dokploy/servers/{server_id}/connection-test", dokployHandler.TestConnection)
	private.HandleFunc("GET /api/v1/integrations/dokploy/servers/{server_id}/applications", dokployHandler.ListApplications)

	protected := config.Admin(private)
	if protected == nil {
		return nil, ErrAdminMiddlewareUnavailable
	}
	root := http.NewServeMux()
	root.HandleFunc("GET /api/v1/integrations/github/app-manifest/callback", githubHandler.CompleteManifest)
	root.HandleFunc("GET /api/v1/integrations/github/app-installations/callback", githubHandler.CompleteInstallation)
	root.Handle("/api/v1/integrations/", protected)
	return root, nil
}
