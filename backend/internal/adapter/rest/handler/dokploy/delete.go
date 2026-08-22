package dokploy

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("server_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	if err := h.servers.Delete(r.Context(), id); err != nil {
		response.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
