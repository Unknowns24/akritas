package mapper

import (
	githubdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/github"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func GitHubManifestRegistrationResultToDTO(value portsin.GitHubManifestRegistrationResult) githubdto.GitHubManifestRegistrationDTO {
	return githubdto.GitHubManifestRegistrationDTO{RegistrationID: value.RegistrationID.String(), FormAction: value.FormAction, Manifest: value.Manifest, State: value.State, ExpiresAt: value.ExpiresAt}
}
