package router

import (
	systemhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/system"
	"github.com/go-chi/chi/v5"
)

func registerSystemPublicRoutes(router chi.Router, handler *systemhandler.Handler) {
	router.Get("/health", handler.Health)
	router.Get("/readiness", handler.Readiness)
}

func registerSystemRoutes(router chi.Router, handler *systemhandler.Handler) {
	router.Get("/system/status", handler.Status)
	router.Post("/system/diagnostics", handler.Diagnostics)
}
