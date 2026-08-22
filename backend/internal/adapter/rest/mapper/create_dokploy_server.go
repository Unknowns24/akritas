package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func CreateDokployServerToCommand(value dto.CreateDokployServerRequestDTO) portsin.CreateDokployServerCommand {
	return portsin.CreateDokployServerCommand{Name: value.Name, BaseURL: value.BaseURL, APICredential: value.APICredential}
}
