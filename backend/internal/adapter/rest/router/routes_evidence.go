package router

import (
	evidencehandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/evidence"
	"github.com/go-chi/chi/v5"
)

func registerEvidenceRoutes(router chi.Router, handler *evidencehandler.Handler) {
	router.Get("/investigations/{investigation_id}/evidence", handler.List)
}
