package auth

import (
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

// Logout assumes RequireSession already ran: a session is always present
// in the request context by the time this handler runs.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.SessionFromContext(r.Context())

	if err := h.logoutAdministrator.Execute(r.Context(), session); err != nil {
		writeSessionError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     response.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.sessionCookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	w.WriteHeader(http.StatusNoContent)
}
