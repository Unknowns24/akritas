package router

import (
	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/go-chi/chi/v5"
)

func registerAuthRoutes(
	router chi.Router,
	handler *authhandler.Handler,
	authenticate portsin.AuthenticateSessionUseCase,
	allowedOrigins []string,
) {
	router.Route("/auth", func(auth chi.Router) {
		auth.Get("/setup-status", handler.GetSetupStatus)
		auth.Post("/setup", handler.StartAdministratorSetup)
		auth.Post("/setup/verify", handler.VerifyAdministratorSetup)
		auth.Post("/login", handler.Login)

		requireSession := middleware.RequireSession(authenticate)
		auth.With(requireSession).Get("/session", handler.GetCurrentSession)
		auth.With(requireSession, middleware.RequireAllowedOrigin(allowedOrigins)).Delete("/session", handler.Logout)
	})
}
