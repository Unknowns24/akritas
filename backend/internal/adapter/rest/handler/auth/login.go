package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	authusecase "github.com/Unknowns24/akritas/backend/internal/usecase/auth"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authdto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

func writeLoginError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
		return
	}
	if errors.Is(err, authusecase.ErrLoginRateLimited) {
		w.Header().Set("Retry-After", "60")
		response.WriteError(w, r, http.StatusTooManyRequests, response.CodeRateLimited,
			"login rate limit exceeded", "Alcanzaste el límite de intentos. Probá nuevamente más tarde.")
		return
	}
	response.WriteInternalError(w, r)
}
