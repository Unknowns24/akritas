package project

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/envelope"
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/httperr"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/utils"
	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetMonitoring(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		httperr.Write(w, apperr.ErrInvalidProjectCommand.Wrap(err), middleware.RequestID(r))
		return
	}
	configuration, err := h.monitoringGetter.GetMonitoring(r.Context(), id)
	if err != nil {
		httperr.Write(w, err, middleware.RequestID(r))
		return
	}
	utils.WriteJSON(w, http.StatusOK, envelope.DataEnvelopeDTO[projectdto.MonitoringConfigurationDTO]{Data: projectdto.FromMonitoring(configuration)})
}
