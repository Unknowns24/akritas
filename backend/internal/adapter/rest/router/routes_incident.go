package router

import (
	incidenthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/incident"
	"github.com/go-chi/chi/v5"
)

func registerIncidentRoutes(router chi.Router, handler *incidenthandler.Handler) {
	router.Get("/incidents", handler.List)
	router.Get("/incidents/{incident_id}", handler.Get)
	router.Get("/incidents/{incident_id}/timeline", handler.ListTimeline)
	router.Get("/incidents/{incident_id}/log-events", handler.ListLogEvents)
}
