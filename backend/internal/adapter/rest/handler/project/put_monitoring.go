package project

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) PutMonitoring(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("project_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	var body projectdto.MonitoringConfigurationRequestDTO
	if request.DecodeJSON(w, r, &body) != nil {
		response.Invalid(w, r)
		return
	}
	configuration, err := mapper.MonitoringConfigurationToDomain(&body)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	updated, err := h.projects.PutMonitoring(r.Context(), id, configuration)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[projectdto.MonitoringConfigurationDTO]{Data: mapper.MonitoringConfigurationToDTO(updated)})
}
