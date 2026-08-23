package mapper

import (
	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func SetupStatus(output in.SetupStatus) authdto.SetupStatusResponseDTO {
	return authdto.SetupStatusResponseDTO{Data: authdto.SetupStatusDTO{
		SetupRequired: output.SetupRequired, RegistrationOpen: output.RegistrationOpen,
	}}
}
