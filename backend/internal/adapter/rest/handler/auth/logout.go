package auth

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
)

// Logout assumes RequireSession already ran: a session is always present
// in the request context by the time this handler runs.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.SessionFromContext(r.Context())

	if err := h.logoutAdministrator.Execute(r.Context(), session); err != nil {
		writeSessionError(w, r, err)
		return
	}

	expireSessionCookie(w, h.sessionCookieSecure)
	w.WriteHeader(http.StatusNoContent)
}
