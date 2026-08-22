package auth

import (
	"net/http"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) GetSetupStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.getSetupStatus.Execute(r.Context())
	if err != nil {
		response.WriteInternalError(w, r)
		return
	}
	response.WriteJSON(w, http.StatusOK, authdto.SetupStatusResponse{
		Data: authdto.SetupStatus{SetupRequired: status.SetupRequired, RegistrationOpen: status.RegistrationOpen},
	})
}
