package remediation

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	operationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/operation"
	remediationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/remediation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) GetIncidentRemediation(w http.ResponseWriter, r *http.Request) {
	incidentID, err := uuid.Parse(r.PathValue("incident_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	value, err := h.remediations.GetIncidentRemediation(r.Context(), incidentID)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[remediationdto.RemediationDTO]{Data: mapper.RemediationToDTO(*value)})
}

func (h *Handler) StartIncidentRemediation(w http.ResponseWriter, r *http.Request) {
	incidentID, err := uuid.Parse(r.PathValue("incident_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	idempotencyKey, err := request.IdempotencyKey(r)
	if err != nil {
		response.Invalid(w, r)
		return
	}
	operation, err := h.remediations.StartIncidentRemediation(r.Context(), portsin.StartIncidentRemediationCommand{
		IncidentID: incidentID, IdempotencyKey: idempotencyKey, WorkspaceRoot: h.workspaceRoot,
	})
	if err != nil {
		response.Error(w, r, err)
		return
	}
	w.Header().Set("Retry-After", "5")
	response.JSON(w, http.StatusAccepted, commondto.DataResponseDTO[operationdto.OperationDTO]{Data: mapper.OperationToDTO(*operation)})
}

func (h *Handler) GetRemediation(w http.ResponseWriter, r *http.Request) {
	remediationID, err := uuid.Parse(r.PathValue("remediation_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	value, err := h.remediations.GetRemediation(r.Context(), remediationID)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[remediationdto.RemediationDTO]{Data: mapper.RemediationToDTO(*value)})
}

func (h *Handler) ListValidationResults(w http.ResponseWriter, r *http.Request) {
	remediationID, err := uuid.Parse(r.PathValue("remediation_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedSorts: map[string]struct{}{"started_at": {}, "finished_at": {}, "id": {}},
		DefaultSort:  []ukerpagination.SortExpression{{Field: "started_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil {
		response.Invalid(w, r)
		return
	}
	page, err := h.remediations.ListValidationResults(r.Context(), remediationID, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]remediationdto.ValidationResultDTO, 0, len(page.Items))
	for _, value := range page.Items {
		items = append(items, mapper.ValidationResultToDTO(value))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[remediationdto.ValidationResultDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[remediationdto.ValidationResultDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func (h *Handler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	remediationID, err := uuid.Parse(r.PathValue("remediation_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	idempotencyKey, err := request.IdempotencyKey(r)
	if err != nil {
		response.Invalid(w, r)
		return
	}
	operation, err := h.remediations.QueueRemediationPullRequest(r.Context(), portsin.CreateRemediationPullRequestCommand{
		RemediationID: remediationID, IdempotencyKey: idempotencyKey, WorkspaceRoot: h.workspaceRoot,
	})
	if err != nil {
		response.Error(w, r, err)
		return
	}
	w.Header().Set("Retry-After", "5")
	response.JSON(w, http.StatusAccepted, commondto.DataResponseDTO[operationdto.OperationDTO]{Data: mapper.OperationToDTO(*operation)})
}
