package investigation

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	investigationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/investigation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("investigation_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	value, err := h.investigations.GetInvestigation(r.Context(), id)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[investigationdto.InvestigationDTO]{Data: mapper.InvestigationToDTO(*value)})
}
