package mapper

import (
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func DokployServerToDTO(value domain.DokployServer) dto.DokployServerDTO {
	result := dto.DokployServerDTO{
		ID: value.ID.String(), Name: value.Name, BaseURL: value.BaseURL, ServerIdentifier: value.ServerIdentifier,
		ConnectionStatus: string(value.ConnectionStatus), CredentialConfigured: value.CredentialConfigured,
		ApplicationCount: value.ApplicationCount, CreatedAt: value.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: value.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if value.LastSyncedAt != nil {
		formatted := value.LastSyncedAt.UTC().Format(time.RFC3339)
		result.LastSyncedAt = &formatted
	}
	return result
}
