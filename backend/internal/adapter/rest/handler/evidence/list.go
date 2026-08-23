package evidence

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	evidencedto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/evidence"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	restpagination "github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	investigationID, err := uuid.Parse(r.PathValue("investigation_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	params, err := restpagination.Parse(r, h.paging, restpagination.Policy{
		AllowedFilters: map[string]struct{}{"type_in": {}},
		AllowedSorts:   map[string]struct{}{"created_at": {}, "id": {}},
		DefaultSort:    []ukerpagination.SortExpression{{Field: "created_at", Direction: ukerpagination.DirectionDesc}, {Field: "id", Direction: ukerpagination.DirectionDesc}},
	})
	if err != nil || !validEvidenceTypes(params.Filters["type_in"]) {
		response.Invalid(w, r)
		return
	}
	page, err := h.evidence.ListInvestigationEvidence(r.Context(), investigationID, params)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	items := make([]evidencedto.EvidenceDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, mapper.EvidenceToDTO(item))
	}
	built, err := restpagination.BuildPage(params, paging.Slice[evidencedto.EvidenceDTO]{Items: items, Total: page.Total}, h.paging, nil)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.ListResponseDTO[evidencedto.EvidenceDTO]{Data: built.Data, Paging: mapper.PagingToDTO(built.Paging)})
}

func validEvidenceTypes(value string) bool {
	if value == "" {
		return true
	}
	for _, item := range strings.Split(value, ",") {
		if domain.EvidenceType(strings.TrimSpace(item)).Validate() != nil {
			return false
		}
	}
	return true
}
