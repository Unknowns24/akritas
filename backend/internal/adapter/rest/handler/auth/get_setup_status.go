package auth

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) GetSetupStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.getSetupStatus.Execute(r.Context())
	if err != nil {
		response.WriteInternalError(w, r)
		return
	}
	response.WriteJSON(w, http.StatusOK, mapper.SetupStatus(status))
}
