package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func UpdateDokployServerToCommand(value dto.UpdateDokployServerRequestDTO) portsin.UpdateDokployServerCommand {
	return portsin.UpdateDokployServerCommand{Name: value.Name, BaseURL: value.BaseURL, APICredential: value.APICredential}
}
