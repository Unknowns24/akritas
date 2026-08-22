package main

import (
	"log"
	"net/http"

	"github.com/Unknowns24/akritas/backend/config"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/dokployserver"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	projectrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/external/integration"
	projecthandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	projectuc "github.com/Unknowns24/akritas/backend/internal/usecase/project"
)

func main() {
	cfg := config.Load()
	db, err := postgres.Open(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := migrations.Run(db); err != nil {
		log.Fatal(err)
	}
	useCase := projectuc.NewUseCase(
		projectrepo.NewRepository(db),
		githubaccount.NewRepository(db),
		dokployserver.NewRepository(db),
		integration.NewSnapshotResolver(),
	)
	handler := projecthandler.NewHandler(useCase, useCase, useCase, useCase, useCase, useCase, cfg.PaginationSecret)
	server := router.New(handler, middleware.RejectAllSessions{})
	log.Printf("akritas listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, server); err != nil {
		log.Fatal(err)
	}
}
