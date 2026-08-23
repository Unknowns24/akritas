package auth

import (
	"errors"
	"net/http"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authdto.LoginRequestDTO
	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest,
			"malformed request body", "La solicitud contiene datos inválidos.")
		return
	}

	if details := validateLoginRequest(req); len(details) > 0 {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest,
			"malformed request body", "La solicitud contiene datos inválidos.", details...)
		return
	}

	output, err := h.loginAdministrator.Execute(r.Context(), in.LoginAdministratorInput{
		Email:        req.Email,
		Password:     req.Password,
		TOTPCode:     req.TotpCode,
		RateLimitKey: clientIP(r),
	})
	if err != nil {
		writeLoginError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     response.SessionCookieName,
		Value:    output.SessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.sessionCookieSecure,
		SameSite: http.SameSiteLaxMode,
		Expires:  output.AbsoluteExpiresAt,
	})
	w.Header().Set("Cache-Control", "no-store")
	response.WriteJSON(w, http.StatusOK, mapper.AdministratorSession(
		output.Administrator, output.AuthenticatedAt, output.IdleExpiresAt, output.AbsoluteExpiresAt,
	))
}

func writeLoginError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
		return
	}
	response.WriteInternalError(w, r)
}
