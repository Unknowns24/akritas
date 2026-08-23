package router

import (
	remediationhandler "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/remediation"
	"github.com/go-chi/chi/v5"
)

func registerRemediationRoutes(router chi.Router, handler *remediationhandler.Handler) {
	router.Get("/incidents/{incident_id}/remediation", handler.GetIncidentRemediation)
	router.Post("/incidents/{incident_id}/remediation", handler.StartIncidentRemediation)
	router.Get("/remediations/{remediation_id}", handler.GetRemediation)
	router.Get("/remediations/{remediation_id}/validation-results", handler.ListValidationResults)
	router.Post("/remediations/{remediation_id}/pull-request", handler.CreatePullRequest)
	router.Get("/pull-requests", handler.ListPullRequests)
	router.Get("/pull-requests/{pull_request_id}", handler.GetPullRequest)
}
