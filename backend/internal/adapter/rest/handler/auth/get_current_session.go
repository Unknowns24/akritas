package auth

import (
	"errors"
	"net/http"
	"time"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
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

	response.WriteJSON(w, http.StatusOK, authdto.SessionResponse{
		Data: authdto.Session{
			Administrator: authdto.Administrator{
				ID:          output.Administrator.ID.String(),
				Email:       output.Administrator.Email,
				DisplayName: output.Administrator.DisplayName,
				CreatedAt:   output.Administrator.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   output.Administrator.UpdatedAt.Format(time.RFC3339),
			},
			AuthenticatedAt:   output.AuthenticatedAt.Format(time.RFC3339),
			IdleExpiresAt:     output.IdleExpiresAt.Format(time.RFC3339),
			AbsoluteExpiresAt: output.AbsoluteExpiresAt.Format(time.RFC3339),
		},
	})
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
