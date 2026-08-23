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

func (h *Handler) ListTimeline(w http.ResponseWriter, r *http.Request) {
	incidentID, err := uuid.Parse(r.PathValue("incident_id"))
	if err != nil {
		response.Error(w, r, domain.ErrIncidentNotFound)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{},
		AllowedSorts:   map[string]struct{}{"occurred_at": {}, "id": {}},
		DefaultSort:    []ukerpagination.SortExpression{{Field: "occurred_at", Direction: ukerpagination.DirectionAsc}, {Field: "id", Direction: ukerpagination.DirectionAsc}},
	})
	if err != nil {
		response.Invalid(w, r)
		return
	}
	page, err := h.incidents.ListTimeline(r.Context(), incidentID, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]incidentdto.TimelineEventDTO, 0, len(page.Items))
	for _, value := range page.Items {
		items = append(items, mapper.TimelineEventToDTO(value))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[incidentdto.TimelineEventDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[incidentdto.TimelineEventDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}
