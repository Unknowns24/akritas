package project

import "github.com/google/uuid"

type UpdateProjectRequestDTO struct {
	Name                 *string                  `json:"name"`
	Description          *string                  `json:"description"`
	GitHubAccountID      *uuid.UUID               `json:"github_account_id"`
	RepositoryIdentifier *string                  `json:"repository_identifier"`
	DefaultBranch        *string                  `json:"default_branch"`
	DokploySource        *DokploySourceRequestDTO `json:"dokploy_source"`
}
