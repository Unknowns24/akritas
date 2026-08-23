package router

import (
	automationhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/automation"
	"github.com/go-chi/chi/v5"
)

func registerAutomationRoutes(router chi.Router, handler *automationhandler.Handler) {
	router.Get("/settings/automation", handler.GetPolicy)
	router.Put("/settings/automation", handler.PutPolicy)
}
