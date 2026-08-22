package router

import (
	"net/http"

	projecthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/system"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func New(projects *projecthandler.Handler, sessions middleware.SessionValidator) http.Handler {
	root := chi.NewRouter()
	root.Use(middleware.RequestIDMiddleware)
	root.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", system.Health)
		r.Get("/readiness", system.Readiness)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireSession(sessions))
			r.Get("/projects", projects.List)
			r.Post("/projects", projects.Create)
			r.Get("/projects/{project_id}", projects.Get)
			r.Patch("/projects/{project_id}", projects.Update)
			r.Get("/projects/{project_id}/monitoring-configuration", projects.GetMonitoring)
			r.Put("/projects/{project_id}/monitoring-configuration", projects.PutMonitoring)
		})
	})
	return root
}
