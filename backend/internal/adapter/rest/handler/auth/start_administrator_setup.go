package auth

import (
	"errors"
	"net"
	"net/http"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func (h *Handler) StartAdministratorSetup(w http.ResponseWriter, r *http.Request) {
	var req authdto.SetupRequestDTO
	if err := request.DecodeJSON(w, r, &req); err != nil {
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
	response.WriteJSON(w, http.StatusCreated, mapper.TOTPEnrollment(output))
}

func writeStartAdministratorSetupError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, domain.ErrInvalidBootstrapToken) {
		response.Invalid(w, r)
		return
	}
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		response.WriteDomainError(w, r, domainErr)
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
