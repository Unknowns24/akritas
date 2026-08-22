package auth

import (
	"net/mail"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
)

// validateSetupRequest checks the transport-level shape declared by
// SetupRequest in the OpenAPI contract (lengths/format). It intentionally
// does not distinguish auth-relevant failures (ADR-008 keeps those generic);
// it only rejects structurally malformed input.
func validateSetupRequest(req authdto.SetupRequest) []dto.ErrorDetail {
	var details []dto.ErrorDetail

	if parsed, err := mail.ParseAddress(req.Email); err != nil || parsed.Address != req.Email || len(req.Email) > 254 {
		details = append(details, dto.ErrorDetail{Field: "email", Reason: "Debe ser un email válido de hasta 254 caracteres."})
	}
	if l := len(req.DisplayName); l < 1 || l > 100 {
		details = append(details, dto.ErrorDetail{Field: "display_name", Reason: "Debe tener entre 1 y 100 caracteres."})
	}
	if l := len(req.Password); l < 12 || l > 128 {
		details = append(details, dto.ErrorDetail{Field: "password", Reason: "Debe tener entre 12 y 128 caracteres."})
	}
	if l := len(req.BootstrapToken); l < 32 || l > 512 {
		details = append(details, dto.ErrorDetail{Field: "bootstrap_token", Reason: "Debe tener entre 32 y 512 caracteres."})
	}

	return details
}
