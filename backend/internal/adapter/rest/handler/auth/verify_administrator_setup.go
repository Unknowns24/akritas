package auth

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

var totpCodePattern = regexp.MustCompile(`^[0-9]{6}$`)

func (h *Handler) VerifyAdministratorSetup(w http.ResponseWriter, r *http.Request) {
	var req authdto.TOTPEnrollmentVerificationRequestDTO
	if err := request.DecodeJSON(w, r, &req); err != nil {
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

func validateVerifyRequest(req authdto.TOTPEnrollmentVerificationRequestDTO) []commondto.ErrorDetailDTO {
	var details []commondto.ErrorDetailDTO
	if _, err := uuid.Parse(req.EnrollmentID); err != nil {
		details = append(details, commondto.ErrorDetailDTO{Field: "enrollment_id", Reason: "Debe ser un UUID válido."})
	}
	if !totpCodePattern.MatchString(req.TotpCode) {
		details = append(details, commondto.ErrorDetailDTO{Field: "totp_code", Reason: "Debe ser un código de 6 dígitos."})
	}
	return details
}

func writeVerifyAdministratorSetupError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
		return
	}
	response.WriteInternalError(w, r)
}
