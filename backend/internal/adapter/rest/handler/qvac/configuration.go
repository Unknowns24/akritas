package qvac

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	qvacdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/qvac"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	value, err := h.qvac.GetConfiguration(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[qvacdto.ConfigurationDTO]{Data: mapper.QvacConfigurationToDTO(value)})
}

func (h *Handler) PutConfiguration(w http.ResponseWriter, r *http.Request) {
	var body qvacdto.PutConfigurationRequestDTO
	if err := request.DecodeJSON(w, r, &body); err != nil {
		response.Invalid(w, r)
		return
	}
	value, err := h.qvac.PutConfiguration(r.Context(), mapper.PutQvacConfigurationToCommand(body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[qvacdto.ConfigurationDTO]{Data: mapper.QvacConfigurationToDTO(value)})
}
