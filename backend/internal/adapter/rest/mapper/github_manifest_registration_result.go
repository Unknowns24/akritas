package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func GitHubManifestRegistrationResultToDTO(value portsin.GitHubManifestRegistrationResult) dto.GitHubManifestRegistrationDTO {
	return dto.GitHubManifestRegistrationDTO{RegistrationID: value.RegistrationID.String(), FormAction: value.FormAction, Manifest: value.Manifest, State: value.State, ExpiresAt: value.ExpiresAt}
}
