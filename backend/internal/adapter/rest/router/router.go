package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
)

type Dependencies struct {
	Auth *authhandler.Handler
}

func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Get("/setup-status", deps.Auth.GetSetupStatus)
		r.Post("/setup", deps.Auth.StartAdministratorSetup)
		r.Post("/setup/verify", deps.Auth.VerifyAdministratorSetup)
	})

	return r
}
