package mapper

import (
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func UpdateDokployServerToCommand(value dokploydto.UpdateDokployServerRequestDTO) portsin.UpdateDokployServerCommand {
	return portsin.UpdateDokployServerCommand{Name: value.Name, BaseURL: value.BaseURL, APICredential: value.APICredential}
}
