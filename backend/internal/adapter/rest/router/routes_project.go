package router

import (
	projecthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/project"
	"github.com/go-chi/chi/v5"
)

func registerProjectRoutes(router chi.Router, handler *projecthandler.Handler) {
	router.Get("/projects", handler.List)
	router.Post("/projects", handler.Create)
	router.Get("/projects/{project_id}", handler.Get)
	router.Patch("/projects/{project_id}", handler.Update)
	router.Delete("/projects/{project_id}", handler.Delete)
	router.Get("/projects/{project_id}/monitoring-configuration", handler.GetMonitoring)
	router.Put("/projects/{project_id}/monitoring-configuration", handler.PutMonitoring)
}
