package dokploy

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

func (h *Handler) ListComposes(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("server_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{"name_like": {}}, Scope: "dokploy_compose:" + id.String(), BoundaryField: "provider_offset",
	})
	if err != nil || len(params.Filters["name_like"]) > 200 {
		response.Invalid(w, r)
		return
	}
	page, err := h.servers.ListComposes(r.Context(), id, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]dokploydto.DokployComposeDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.DokployComposeToDTO(item))
	}
	built, err := restpagination.BuildProviderPage(params, paging.Slice[dokploydto.DokployComposeDTO]{Items: items, Total: page.Total, NextBoundary: page.NextBoundary, PrevBoundary: page.PrevBoundary}, h.paging)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[dokploydto.DokployComposeDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func (h *Handler) ListComposeServices(w http.ResponseWriter, r *http.Request) {
	serverID, err := uuid.Parse(r.PathValue("server_id"))
	composeID := strings.TrimSpace(r.PathValue("compose_id"))
	if err != nil || composeID == "" || len(composeID) > 255 {
		response.Invalid(w, r)
		return
	}
	query := r.URL.Query()
	refresh := false
	if query.Has("refresh") {
		if len(query["refresh"]) != 1 {
			response.Invalid(w, r)
			return
		}
		raw := query.Get("refresh")
		if raw != "true" && raw != "false" {
			response.Invalid(w, r)
			return
		}
		refresh = raw == "true"
	}
	if len(query) > 1 || (len(query) == 1 && !query.Has("refresh")) {
		response.Invalid(w, r)
		return
	}
	services, err := h.servers.ListComposeServices(r.Context(), serverID, composeID, refresh)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]dokploydto.DokployComposeServiceDTO, 0, len(services))
	for _, service := range services {
		items = append(items, mapper.DokployComposeServiceToDTO(service))
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[[]dokploydto.DokployComposeServiceDTO]{Data: items})
}
