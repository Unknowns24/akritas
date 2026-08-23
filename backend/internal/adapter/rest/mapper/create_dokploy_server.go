package mapper

import (
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func CreateDokployServerToCommand(value dokploydto.CreateDokployServerRequestDTO) portsin.CreateDokployServerCommand {
	return portsin.CreateDokployServerCommand{Name: value.Name, BaseURL: value.BaseURL, APICredential: value.APICredential}
}
