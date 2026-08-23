package router

import (
	githubhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/github"
	"github.com/go-chi/chi/v5"
)

func registerGitHubCallbackRoutes(router chi.Router, handler *githubhandler.Handler) {
	router.Get("/github/app-manifest/callback", handler.CompleteManifest)
	router.Get("/github/app-installations/callback", handler.CompleteInstallation)
}

func registerGitHubRoutes(router chi.Router, handler *githubhandler.Handler) {
	router.Route("/github", func(github chi.Router) {
		github.Get("/accounts", handler.ListAccounts)
		github.Post("/accounts", handler.CreateAccount)
		github.Get("/accounts/{account_id}", handler.GetAccount)
		github.Patch("/accounts/{account_id}", handler.UpdateAccount)
		github.Delete("/accounts/{account_id}", handler.DeleteAccount)
		github.Post("/accounts/{account_id}/connection-test", handler.TestConnection)
		github.Get("/accounts/{account_id}/repositories", handler.ListRepositories)
		github.Post("/app-manifest/registrations", handler.StartManifest)
	})
}
