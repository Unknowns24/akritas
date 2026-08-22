package projectdto

import (
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
	"github.com/google/uuid"
)

type UpdateProjectRequestDTO struct {
	Name                  *string    `json:"name"`
	Description           *string    `json:"description"`
	GitHubAccountID       *uuid.UUID `json:"github_account_id"`
	RepositoryIdentifier  *string    `json:"repository_identifier"`
	DefaultBranch         *string    `json:"default_branch"`
	DokployServerID       *uuid.UUID `json:"dokploy_server_id"`
	ApplicationIdentifier *string    `json:"application_identifier"`
}

func (r UpdateProjectRequestDTO) Command(id uuid.UUID) inproject.UpdateCommand {
	return inproject.UpdateCommand{
		ID: id, Name: r.Name, Description: r.Description, GitHubAccountID: r.GitHubAccountID,
		RepositoryIdentifier: r.RepositoryIdentifier, DefaultBranch: r.DefaultBranch,
		DokployServerID: r.DokployServerID, ApplicationIdentifier: r.ApplicationIdentifier,
	}
}
