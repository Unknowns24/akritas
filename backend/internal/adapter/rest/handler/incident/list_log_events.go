package incident

import (
	"net/http"

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

func (h *Handler) ListLogEvents(w http.ResponseWriter, r *http.Request) {
	incidentID, err := uuid.Parse(r.PathValue("incident_id"))
	if err != nil {
		response.Error(w, r, domain.ErrIncidentNotFound)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{AllowedFilters: map[string]struct{}{}, AllowedSorts: map[string]struct{}{"timestamp": {}, "id": {}}, DefaultSort: []ukerpagination.SortExpression{{Field: "timestamp", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}}})
	if err != nil {
		response.Invalid(w, r)
		return
	}
	page, err := h.incidents.ListLogEvents(r.Context(), incidentID, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]incidentdto.LogEventDTO, 0, len(page.Items))
	for _, value := range page.Items {
		items = append(items, mapper.LogEventToDTO(value))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[incidentdto.LogEventDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[incidentdto.LogEventDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}
