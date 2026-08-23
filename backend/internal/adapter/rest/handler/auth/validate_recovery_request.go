package auth

import (
	"net/mail"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
)

func validateRecoveryRequest(req authdto.RecoveryRequestDTO) []commondto.ErrorDetailDTO {
	var details []commondto.ErrorDetailDTO
	if parsed, err := mail.ParseAddress(req.Email); err != nil || parsed.Address != req.Email || len(req.Email) > 254 {
		details = append(details, commondto.ErrorDetailDTO{Field: "email", Reason: "Debe ser un email válido de hasta 254 caracteres."})
	}
	if l := len(req.NewPassword); l < 12 || l > 128 {
		details = append(details, commondto.ErrorDetailDTO{Field: "new_password", Reason: "Debe tener entre 12 y 128 caracteres."})
	}
	if l := len(req.BootstrapToken); l < 32 || l > 512 {
		details = append(details, commondto.ErrorDetailDTO{Field: "bootstrap_token", Reason: "Debe tener entre 32 y 512 caracteres."})
	}
	return details
}
