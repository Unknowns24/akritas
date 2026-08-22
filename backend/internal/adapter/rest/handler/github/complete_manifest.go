package github

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) CompleteManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	result, err := h.apps.CompleteManifest(r.Context(), r.URL.Query().Get("code"), r.URL.Query().Get("state"))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	w.Header().Set("Location", result.RedirectURL)
	w.WriteHeader(http.StatusSeeOther)
}
