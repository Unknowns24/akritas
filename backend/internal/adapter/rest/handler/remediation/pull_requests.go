package remediation

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	remediationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/remediation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) ListPullRequests(w http.ResponseWriter, r *http.Request) {
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{"project_id_eq": {}, "incident_id_eq": {}},
		AllowedSorts:   map[string]struct{}{"created_at": {}, "id": {}},
		DefaultSort:    []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil || !validPullRequestFilters(params.Filters) {
		response.Invalid(w, r)
		return
	}
	page, err := h.remediations.ListPullRequests(r.Context(), params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]remediationdto.PullRequestDTO, 0, len(page.Items))
	for _, value := range page.Items {
		items = append(items, mapper.PullRequestProjectionToDTO(value))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[remediationdto.PullRequestDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[remediationdto.PullRequestDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func (h *Handler) GetPullRequest(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("pull_request_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	value, err := h.remediations.GetPullRequest(r.Context(), id)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[remediationdto.PullRequestDTO]{Data: mapper.PullRequestProjectionToDTO(*value)})
}

func validPullRequestFilters(filters map[string]string) bool {
	for _, key := range []string{"project_id_eq", "incident_id_eq"} {
		if value := filters[key]; value != "" {
			if _, err := uuid.Parse(value); err != nil {
				return false
			}
		}
	}
	return true
}
