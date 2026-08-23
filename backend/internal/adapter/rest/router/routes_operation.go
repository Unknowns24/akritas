package router

import (
	operationhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/operation"
	"github.com/go-chi/chi/v5"
)

func registerOperationRoutes(router chi.Router, handler *operationhandler.Handler) {
	router.Get("/operations/{operation_id}", handler.Get)
}
