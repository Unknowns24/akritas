package router

import (
	"errors"
	"net/http"

	resthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler"
	authmiddleware "github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	ErrAdminMiddlewareUnavailable = errors.New("administrator middleware is unavailable")
	ErrInvalidRouterConfiguration = errors.New("invalid integration router configuration")
)

type AdminMiddleware func(http.Handler) http.Handler

type Config struct {
	Handlers       *resthandler.Handlers
	Admin          AdminMiddleware
	Authenticate   portsin.AuthenticateSessionUseCase
	AllowedOrigins []string
}

func New(config Config) (http.Handler, error) {
	if config.Admin == nil {
		return nil, ErrAdminMiddlewareUnavailable
	}

	if config.Handlers == nil ||
		config.Handlers.AuthHandler == nil ||
		config.Handlers.GitHubHandler == nil ||
		config.Handlers.DokployHandler == nil ||
		config.Handlers.ProjectHandler == nil ||
		config.Handlers.InvestigationHandler == nil ||
		config.Handlers.OperationHandler == nil ||
		config.Handlers.EvidenceHandler == nil ||
		config.Handlers.IncidentHandler == nil ||
		config.Authenticate == nil ||
		len(config.AllowedOrigins) == 0 {
		return nil, ErrInvalidRouterConfiguration
	}

	if config.Admin(http.NotFoundHandler()) == nil {
		return nil, ErrAdminMiddlewareUnavailable
	}

	root := chi.NewRouter()
	root.Use(chimiddleware.RequestID)
	root.Use(authmiddleware.RecoverPanics)
	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Idempotency-Key"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	root.Use(chimiddleware.GetHead)
	root.Route("/api/v1", func(api chi.Router) {
		registerAuthRoutes(api, config.Handlers.AuthHandler, config.Authenticate, config.AllowedOrigins)
		api.Route("/integrations", func(integrations chi.Router) {
			registerGitHubCallbackRoutes(integrations, config.Handlers.GitHubHandler)
			integrations.Group(func(private chi.Router) {
				private.Use(config.Admin)
				private.Use(authmiddleware.RequireAllowedOrigin(config.AllowedOrigins))
				registerGitHubRoutes(private, config.Handlers.GitHubHandler)
				registerDokployRoutes(private, config.Handlers.DokployHandler)
			})
		})
		api.Group(func(private chi.Router) {
			private.Use(config.Admin)
			private.Use(authmiddleware.RequireAllowedOrigin(config.AllowedOrigins))
			registerProjectRoutes(private, config.Handlers.ProjectHandler)
			registerInvestigationRoutes(private, config.Handlers.InvestigationHandler)
			registerOperationRoutes(private, config.Handlers.OperationHandler)
			registerEvidenceRoutes(private, config.Handlers.EvidenceHandler)
			registerIncidentRoutes(private, config.Handlers.IncidentHandler)
		})
	})
	return root, nil
}
