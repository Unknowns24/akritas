package mapper

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func GitHubAccountToDTO(value domain.GitHubAccount) dto.GitHubAccountDTO {
	result := dto.GitHubAccountDTO{
		ID: value.ID.String(), DisplayName: value.DisplayName, AccountType: string(value.AccountType), AccountIdentifier: value.AccountIdentifier,
		AuthenticationMethod: string(value.AuthenticationMethod), AuthenticationStatus: string(value.AuthenticationStatus),
		CredentialConfigured: value.CredentialConfigured, RepositoryCount: value.RepositoryCount, ManageURL: value.ManageURL,
		CreatedAt: value.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: value.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if value.LastCheckedAt != nil {
		formatted := value.LastCheckedAt.UTC().Format(time.RFC3339)
		result.LastCheckedAt = &formatted
	}
	return result
}
