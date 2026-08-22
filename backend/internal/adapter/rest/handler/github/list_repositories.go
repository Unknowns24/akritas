package github

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (h *Handler) ListRepositories(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("account_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{"name_like": {}}, Scope: "github_repository:" + id.String(), BoundaryField: "provider_page",
	})
	if err != nil || len(params.Filters["name_like"]) > 200 {
		response.Invalid(w, r)
		return
	}
	page, err := h.accounts.ListRepositories(r.Context(), id, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]dto.GitHubRepositoryDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.GitHubRepositoryToDTO(item))
	}
	built, err := restpagination.BuildProviderPage(params, paging.Slice[dto.GitHubRepositoryDTO]{Items: items, Total: page.Total, NextBoundary: page.NextBoundary, PrevBoundary: page.PrevBoundary}, h.paging)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.ListResponseDTO[dto.GitHubRepositoryDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}
