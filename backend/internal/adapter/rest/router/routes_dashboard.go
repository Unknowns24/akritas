package router

import (
	dashboardhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/dashboard"
	"github.com/go-chi/chi/v5"
)

func registerDashboardRoutes(router chi.Router, handler *dashboardhandler.Handler) {
	router.Get("/overview", handler.Overview)
	router.Get("/activity", handler.Activity)
}
