package operation

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	operationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/operation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("operation_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	value, err := h.operations.GetOperation(r.Context(), id)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[operationdto.OperationDTO]{Data: mapper.OperationToDTO(*value)})
}
