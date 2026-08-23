package dashboard

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	dashboarddto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dashboard"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) Activity(w http.ResponseWriter, r *http.Request) {
	filters := map[string]struct{}{"project_id_eq": {}, "incident_id_eq": {}, "type_in": {}}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: filters,
		AllowedSorts:   map[string]struct{}{"occurred_at": {}, "id": {}},
		DefaultSort:    []ukerpagination.SortExpression{{Field: "occurred_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil || !validActivityFilters(params.Filters) {
		response.Invalid(w, r)
		return
	}
	page, err := h.dashboard.ListActivity(r.Context(), params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]dashboarddto.ActivityEventDTO, 0, len(page.Items))
	for _, value := range page.Items {
		items = append(items, mapper.ActivityEventToDTO(value))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[dashboarddto.ActivityEventDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[dashboarddto.ActivityEventDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func validActivityFilters(filters map[string]string) bool {
	for _, key := range []string{"project_id_eq", "incident_id_eq"} {
		if value := filters[key]; value != "" {
			if _, err := uuid.Parse(value); err != nil {
				return false
			}
		}
	}
	if value := filters["type_in"]; value != "" {
		for _, item := range strings.Split(value, ",") {
			if err := domain.ActivityType(strings.TrimSpace(item)).Validate(); err != nil {
				return false
			}
		}
	}
	return true
}
