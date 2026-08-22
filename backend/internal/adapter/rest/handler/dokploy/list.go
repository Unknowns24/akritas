package dokploy

import (
	"net/http"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{"connection_status_in": {}},
		AllowedSorts:   map[string]struct{}{"created_at": {}, "id": {}},
		DefaultSort:    []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil {
		response.Invalid(w, r)
		return
	}
	if _, valid := integrationStatuses(params.Filters["connection_status_in"]); !valid {
		response.Invalid(w, r)
		return
	}
	page, err := h.servers.List(r.Context(), params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]dto.DokployServerDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.DokployServerToDTO(item))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[dto.DokployServerDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.ListResponseDTO[dto.DokployServerDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func integrationStatuses(value string) ([]domain.IntegrationStatus, bool) {
	if value == "" {
		return nil, true
	}
	parts := strings.Split(value, ",")
	statuses := make([]domain.IntegrationStatus, 0, len(parts))
	for _, part := range parts {
		status := domain.IntegrationStatus(strings.TrimSpace(part))
		if status.Validate() != nil {
			return nil, false
		}
		statuses = append(statuses, status)
	}
	return statuses, true
}
