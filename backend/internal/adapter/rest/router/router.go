package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	authmiddleware "github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type Dependencies struct {
	Auth         *authhandler.Handler
	Authenticate in.AuthenticateSessionUseCase
}

func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Get("/setup-status", deps.Auth.GetSetupStatus)
		r.Post("/setup", deps.Auth.StartAdministratorSetup)
		r.Post("/setup/verify", deps.Auth.VerifyAdministratorSetup)
		r.Post("/login", deps.Auth.Login)

		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.RequireSession(deps.Authenticate))
			r.Get("/session", deps.Auth.GetCurrentSession)
			r.Delete("/session", deps.Auth.Logout)
		})
	})

	return r
}
