package mapper

import (
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func UpdateProjectToCommand(id uuid.UUID, value projectdto.UpdateProjectRequestDTO) portsin.UpdateProjectCommand {
	return portsin.UpdateProjectCommand{ID: id, Name: value.Name, Description: value.Description, GitHubAccountID: value.GitHubAccountID, RepositoryIdentifier: value.RepositoryIdentifier, DefaultBranch: value.DefaultBranch, DokployServerID: value.DokployServerID, ApplicationIdentifier: value.ApplicationIdentifier}
}
