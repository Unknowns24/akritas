package investigation

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	investigationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/investigation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	incidentID, err := uuid.Parse(r.PathValue("incident_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedSorts: map[string]struct{}{"created_at": {}, "id": {}},
		DefaultSort:  []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil {
		response.Invalid(w, r)
		return
	}
	page, err := h.investigations.ListIncidentInvestigations(r.Context(), incidentID, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]investigationdto.InvestigationDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.InvestigationToDTO(item))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[investigationdto.InvestigationDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[investigationdto.InvestigationDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}
