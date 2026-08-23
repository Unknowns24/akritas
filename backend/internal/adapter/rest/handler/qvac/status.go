package qvac

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	qvacdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/qvac"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) TestConnection(w http.ResponseWriter, r *http.Request) {
	value, err := h.qvac.TestConnection(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[commondto.ConnectionTestDTO]{Data: mapper.ConnectionTestToDTO(value)})
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	value, err := h.qvac.GetStatus(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[qvacdto.RuntimeStatusDTO]{Data: mapper.QvacStatusToDTO(value)})
}
