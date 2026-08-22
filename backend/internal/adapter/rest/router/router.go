package router

import (
	"errors"
	"net/http"

	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/github"
	authmiddleware "github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
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
	Auth           *authhandler.Handler
	Authenticate   portsin.AuthenticateSessionUseCase
	AllowedOrigins []string
}

func New(config Config) (http.Handler, error) {
	if config.Admin == nil {
		return nil, ErrAdminMiddlewareUnavailable
	}
	if config.Auth == nil || config.Authenticate == nil || len(config.AllowedOrigins) == 0 {
		return nil, ErrInvalidRouterConfiguration
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

	originProtected := authmiddleware.RequireAllowedOrigin(config.AllowedOrigins)
	protected := config.Admin(originProtected(private))
	if protected == nil {
		return nil, ErrAdminMiddlewareUnavailable
	}
	root := http.NewServeMux()
	root.HandleFunc("GET /api/v1/auth/setup-status", config.Auth.GetSetupStatus)
	root.HandleFunc("POST /api/v1/auth/setup", config.Auth.StartAdministratorSetup)
	root.HandleFunc("POST /api/v1/auth/setup/verify", config.Auth.VerifyAdministratorSetup)
	root.HandleFunc("POST /api/v1/auth/login", config.Auth.Login)
	authenticatedAuth := authmiddleware.RequireSession(config.Authenticate)
	root.Handle("GET /api/v1/auth/session", authenticatedAuth(http.HandlerFunc(config.Auth.GetCurrentSession)))
	root.Handle("DELETE /api/v1/auth/session", authenticatedAuth(originProtected(http.HandlerFunc(config.Auth.Logout))))
	root.HandleFunc("GET /api/v1/integrations/github/app-manifest/callback", githubHandler.CompleteManifest)
	root.HandleFunc("GET /api/v1/integrations/github/app-installations/callback", githubHandler.CompleteInstallation)
	root.Handle("/api/v1/integrations/", protected)
	return root, nil
}
