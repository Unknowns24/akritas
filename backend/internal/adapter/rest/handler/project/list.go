package project

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{AllowedFilters: map[string]struct{}{"name_like": {}, "monitoring_status_in": {}}, AllowedSorts: map[string]struct{}{"name": {}, "created_at": {}, "id": {}}, DefaultSort: []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}}})
	if err != nil || len(params.Filters["name_like"]) > 200 || !validMonitoringStatuses(params.Filters["monitoring_status_in"]) {
		response.Invalid(w, r)
		return
	}
	page, err := h.projects.List(r.Context(), params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]projectdto.ProjectSummaryDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.ProjectSummaryToDTO(item))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[projectdto.ProjectSummaryDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[projectdto.ProjectSummaryDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func validMonitoringStatuses(value string) bool {
	if value == "" {
		return true
	}
	for _, item := range strings.Split(value, ",") {
		if domain.MonitoringStatus(strings.TrimSpace(item)).Validate() != nil {
			return false
		}
	}
	return true
}
