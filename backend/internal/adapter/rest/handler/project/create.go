package project

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/envelope"
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/httperr"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/utils"
	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request projectdto.CreateProjectRequestDTO
	if err := utils.DecodeJSON(r, &request); err != nil {
		httperr.Write(w, apperr.ErrInvalidProjectCommand.Wrap(err), middleware.RequestID(r))
		return
	}
	command, err := request.Command()
	if err != nil {
		httperr.Write(w, apperr.ErrInvalidProjectCommand.Wrap(err), middleware.RequestID(r))
		return
	}
	result, err := h.projects.Create(r.Context(), command)
	if err != nil {
		httperr.Write(w, err, middleware.RequestID(r))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, envelope.DataEnvelopeDTO[projectdto.ProjectDTO]{Data: projectdto.FromProject(result)})
}
