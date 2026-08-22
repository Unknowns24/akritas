package dokploy

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (h *Handler) ListApplications(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("server_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{"name_like": {}, "environment_eq": {}}, Scope: "dokploy_application:" + id.String(), BoundaryField: "provider_offset",
	})
	if err != nil || len(params.Filters["name_like"]) > 200 || len(params.Filters["environment_eq"]) > 100 {
		response.Invalid(w, r)
		return
	}
	page, err := h.servers.ListApplications(r.Context(), id, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]dokploydto.DokployApplicationDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.DokployApplicationToDTO(item))
	}
	built, err := restpagination.BuildProviderPage(params, paging.Slice[dokploydto.DokployApplicationDTO]{Items: items, Total: page.Total, NextBoundary: page.NextBoundary, PrevBoundary: page.PrevBoundary}, h.paging)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[dokploydto.DokployApplicationDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}
