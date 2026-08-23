package incident

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := map[string]struct{}{"project_id_eq": {}, "phase_in": {}, "severity_in": {}, "root_cause_status_in": {}, "resolution_status_in": {}, "title_like": {}}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{AllowedFilters: filters, AllowedSorts: map[string]struct{}{"last_seen_at": {}, "first_seen_at": {}, "severity": {}, "id": {}}, DefaultSort: []ukerpagination.SortExpression{{Field: "last_seen_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}}})
	if err != nil || !validIncidentFilters(params.Filters) {
		response.Invalid(w, r)
		return
	}
	page, err := h.incidents.List(r.Context(), params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]incidentdto.IncidentDTO, 0, len(page.Items))
	for _, value := range page.Items {
		items = append(items, mapper.IncidentToDTO(value))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[incidentdto.IncidentDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[incidentdto.IncidentDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func validIncidentFilters(filters map[string]string) bool {
	if id := filters["project_id_eq"]; id != "" {
		if _, err := uuid.Parse(id); err != nil {
			return false
		}
	}
	if len(filters["title_like"]) > 200 {
		return false
	}
	return validEnumList(filters["phase_in"], func(value string) bool { return domain.IncidentPhase(value).Validate() == nil }) &&
		validEnumList(filters["severity_in"], func(value string) bool { return domain.Severity(value).Validate() == nil }) &&
		validEnumList(filters["root_cause_status_in"], func(value string) bool { return domain.RootCauseStatus(value).Validate() == nil }) &&
		validEnumList(filters["resolution_status_in"], func(value string) bool { return domain.ResolutionStatus(value).Validate() == nil })
}

func validEnumList(value string, validate func(string) bool) bool {
	if value == "" {
		return true
	}
	for _, item := range strings.Split(value, ",") {
		if !validate(strings.TrimSpace(item)) {
			return false
		}
	}
	return true
}
