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

func (h *Handler) StartAdministratorRecovery(w http.ResponseWriter, r *http.Request) {
	var req authdto.RecoveryRequestDTO
	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.Invalid(w, r)
		return
	}
	if details := validateRecoveryRequest(req); len(details) > 0 {
		response.WriteError(w, r, http.StatusBadRequest, response.CodeMalformedRequest, "malformed request body", "La solicitud contiene datos inválidos.", details...)
		return
	}
	output, err := h.startAdministratorRecovery.Execute(r.Context(), in.StartAdministratorRecoveryInput{Email: req.Email, NewPassword: req.NewPassword, BootstrapToken: req.BootstrapToken, RateLimitKey: clientIP(r)})
	if err != nil {
		writeRecoveryError(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	response.WriteJSON(w, http.StatusCreated, mapper.TOTPRecoveryEnrollment(output))
}

func writeRecoveryError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
		return
	}
	response.WriteInternalError(w, r)
}
