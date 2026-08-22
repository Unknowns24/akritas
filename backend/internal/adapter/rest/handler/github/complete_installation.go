package github

import (
	"net/http"
	"strconv"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) CompleteInstallation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	installationID, err := strconv.ParseInt(r.URL.Query().Get("installation_id"), 10, 64)
	if err != nil || installationID < 1 {
		response.Invalid(w, r)
		return
	}
	result, err := h.apps.CompleteInstallation(r.Context(), installationID, r.URL.Query().Get("state"))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	w.Header().Set("Location", result.RedirectURL)
	w.WriteHeader(http.StatusSeeOther)
}
