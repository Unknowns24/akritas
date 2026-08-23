package router

import (
	qvachandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/qvac"
	"github.com/go-chi/chi/v5"
)

func registerQvacRoutes(router chi.Router, handler *qvachandler.Handler) {
	router.Get("/integrations/qvac/configuration", handler.GetConfiguration)
	router.Put("/integrations/qvac/configuration", handler.PutConfiguration)
	router.Post("/integrations/qvac/connection-test", handler.TestConnection)
	router.Get("/integrations/qvac/status", handler.Status)
}
