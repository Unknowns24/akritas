package auth

import (
	"net/mail"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
)

// validateLoginRequest checks the transport-level shape declared by
// LoginRequest in the OpenAPI contract. It intentionally does not
// distinguish auth-relevant failures (ADR-008 keeps those generic); it
// only rejects structurally malformed input.
func validateLoginRequest(req authdto.LoginRequestDTO) []commondto.ErrorDetailDTO {
	var details []commondto.ErrorDetailDTO

	if parsed, err := mail.ParseAddress(req.Email); err != nil || parsed.Address != req.Email || len(req.Email) > 254 {
		details = append(details, commondto.ErrorDetailDTO{Field: "email", Reason: "Debe ser un email válido de hasta 254 caracteres."})
	}
	if l := len(req.Password); l < 12 || l > 128 {
		details = append(details, commondto.ErrorDetailDTO{Field: "password", Reason: "Debe tener entre 12 y 128 caracteres."})
	}
	if !totpCodePattern.MatchString(req.TotpCode) {
		details = append(details, commondto.ErrorDetailDTO{Field: "totp_code", Reason: "Debe ser un código de 6 dígitos."})
	}

	return details
}
