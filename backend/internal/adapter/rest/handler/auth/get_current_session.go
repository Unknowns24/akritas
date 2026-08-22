package auth

import (
	"errors"
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// GetCurrentSession assumes RequireSession already ran: a session is
// always present in the request context by the time this handler runs.
func (h *Handler) GetCurrentSession(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.SessionFromContext(r.Context())

	output, err := h.getCurrentSession.Execute(r.Context(), session)
	if err != nil {
		writeSessionError(w, r, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, mapper.AdministratorSession(
		output.Administrator, output.AuthenticatedAt, output.IdleExpiresAt, output.AbsoluteExpiresAt,
	))
}

// writeSessionError is shared by GetCurrentSession and Logout: both act on
// a session already validated by RequireSession, so any error left is
// either a domain sentinel (e.g. the administrator vanished) or an opaque
// infrastructure failure.
func writeSessionError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
		return
	}
	response.WriteInternalError(w, r)
}
