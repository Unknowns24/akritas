package router

import (
	dokployhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/dokploy"
	"github.com/go-chi/chi/v5"
)

func registerDokployRoutes(router chi.Router, handler *dokployhandler.Handler) {
	router.Route("/dokploy", func(dokploy chi.Router) {
		dokploy.Get("/servers", handler.List)
		dokploy.Post("/servers", handler.Create)
		dokploy.Get("/servers/{server_id}", handler.Get)
		dokploy.Patch("/servers/{server_id}", handler.Update)
		dokploy.Delete("/servers/{server_id}", handler.Delete)
		dokploy.Post("/servers/{server_id}/connection-test", handler.TestConnection)
		dokploy.Get("/servers/{server_id}/applications", handler.ListApplications)
		dokploy.Get("/servers/{server_id}/composes", handler.ListComposes)
		dokploy.Get("/servers/{server_id}/composes/{compose_id}/services", handler.ListComposeServices)
	})
}
