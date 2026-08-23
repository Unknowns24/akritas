package router

import (
	investigationhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/investigation"
	"github.com/go-chi/chi/v5"
)

func registerInvestigationRoutes(router chi.Router, handler *investigationhandler.Handler) {
	router.Get("/incidents/{incident_id}/investigations", handler.List)
	router.Post("/incidents/{incident_id}/investigations", handler.Start)
	router.Get("/investigations/{investigation_id}", handler.Get)
}
