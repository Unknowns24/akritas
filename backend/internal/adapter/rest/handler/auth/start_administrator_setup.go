package auth

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	authusecase "github.com/Unknowns24/akritas/backend/internal/usecase/auth"
)

func (h *Handler) StartAdministratorSetup(w http.ResponseWriter, r *http.Request) {
	var req authdto.SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest,
			"malformed request body", "La solicitud contiene datos inválidos.")
		return
	}

	if details := validateSetupRequest(req); len(details) > 0 {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest,
			"malformed request body", "La solicitud contiene datos inválidos.", details...)
		return
	}

	output, err := h.startAdministratorSetup.Execute(r.Context(), in.StartAdministratorSetupInput{
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		Password:       req.Password,
		BootstrapToken: req.BootstrapToken,
		RateLimitKey:   clientIP(r),
	})
	if err != nil {
		writeStartAdministratorSetupError(w, r, err)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	response.WriteJSON(w, http.StatusCreated, authdto.TotpEnrollmentResponse{
		Data: authdto.TotpEnrollment{
			EnrollmentID:   output.EnrollmentID.String(),
			OtpauthURI:     output.OtpauthURI,
			ManualEntryKey: output.ManualEntryKey,
			ExpiresAt:      output.ExpiresAt.Format(time.RFC3339),
		},
	})
}

func writeStartAdministratorSetupError(w http.ResponseWriter, r *http.Request, err error) {
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

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
