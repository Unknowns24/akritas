package project

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/envelope"
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/httperr"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/utils"
	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) PutMonitoring(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		httperr.Write(w, apperr.ErrInvalidProjectCommand.Wrap(err), middleware.RequestID(r))
		return
	}
	var request projectdto.MonitoringConfigurationDTO
	if err := utils.DecodeJSON(r, &request); err != nil {
		httperr.Write(w, apperr.ErrInvalidProjectCommand.Wrap(err), middleware.RequestID(r))
		return
	}
	configuration, err := request.Domain()
	if err != nil {
		httperr.Write(w, err, middleware.RequestID(r))
		return
	}
	saved, err := h.monitoringPutter.PutMonitoring(r.Context(), inproject.MonitoringCommand{ProjectID: id, MonitoringConfiguration: configuration})
	if err != nil {
		httperr.Write(w, err, middleware.RequestID(r))
		return
	}
	utils.WriteJSON(w, http.StatusOK, envelope.DataEnvelopeDTO[projectdto.MonitoringConfigurationDTO]{Data: projectdto.FromMonitoring(saved)})
}
