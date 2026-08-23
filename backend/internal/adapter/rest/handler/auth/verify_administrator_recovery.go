package auth

import (
	"net/http"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func (h *Handler) VerifyAdministratorRecovery(w http.ResponseWriter, r *http.Request) {
	var req authdto.TOTPEnrollmentVerificationRequestDTO
	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.Invalid(w, r)
		return
	}
	if details := validateVerifyRequest(req); len(details) > 0 {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest, "malformed request body", "La solicitud contiene datos inválidos.", details...)
		return
	}
	output, err := h.verifyAdministratorRecovery.Execute(r.Context(), in.VerifyAdministratorRecoveryInput{EnrollmentID: req.EnrollmentID, TOTPCode: req.TotpCode, RateLimitKey: clientIP(r)})
	if err != nil {
		writeRecoveryError(w, r, err)
		return
	}
	setSessionCookie(w, h.sessionCookieSecure, h.sessionCookieSameSite, output.SessionToken, output.AuthenticatedAt, output.AbsoluteExpiresAt)
	w.Header().Set("Cache-Control", "no-store")
	response.WriteJSON(w, http.StatusOK, mapper.AdministratorSession(output.Administrator, output.AuthenticatedAt, output.IdleExpiresAt, output.AbsoluteExpiresAt))
}
