package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	authusecase "github.com/Unknowns24/akritas/backend/internal/usecase/auth"
)

var totpCodePattern = regexp.MustCompile(`^[0-9]{6}$`)

func (h *Handler) VerifyAdministratorSetup(w http.ResponseWriter, r *http.Request) {
	var req authdto.TotpEnrollmentVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest,
			"malformed request body", "La solicitud contiene datos inválidos.")
		return
	}

	if details := validateVerifyRequest(req); len(details) > 0 {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest,
			"malformed request body", "La solicitud contiene datos inválidos.", details...)
		return
	}

	output, err := h.verifyAdministratorSetup.Execute(r.Context(), in.VerifyAdministratorSetupInput{
		EnrollmentID: req.EnrollmentID,
		TOTPCode:     req.TotpCode,
		RateLimitKey: clientIP(r),
	})
	if err != nil {
		writeVerifyAdministratorSetupError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
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

func validateVerifyRequest(req authdto.TotpEnrollmentVerificationRequest) []dto.ErrorDetail {
	var details []dto.ErrorDetail
	if _, err := uuid.Parse(req.EnrollmentID); err != nil {
		details = append(details, dto.ErrorDetail{Field: "enrollment_id", Reason: "Debe ser un UUID válido."})
	}
	if !totpCodePattern.MatchString(req.TotpCode) {
		details = append(details, dto.ErrorDetail{Field: "totp_code", Reason: "Debe ser un código de 6 dígitos."})
	}
	return details
}

func writeVerifyAdministratorSetupError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
		return
	}
	if errors.Is(err, authusecase.ErrSetupRateLimited) {
		w.Header().Set("Retry-After", "60")
		response.WriteError(w, r, http.StatusTooManyRequests, response.CodeRateLimited,
			"administrator setup rate limit exceeded", "Alcanzaste el límite de intentos. Probá nuevamente más tarde.")
		return
	}
	response.WriteInternalError(w, r)
}
