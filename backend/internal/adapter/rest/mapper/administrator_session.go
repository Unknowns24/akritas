package mapper

import (
	"time"

	authdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/auth"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func AdministratorSession(
	administrator domain.Administrator,
	authenticatedAt, idleExpiresAt, absoluteExpiresAt time.Time,
) authdto.SessionResponseDTO {
	return authdto.SessionResponseDTO{Data: authdto.SessionDTO{
		Administrator: authdto.AdministratorDTO{
			ID: administrator.ID.String(), Email: administrator.Email, DisplayName: administrator.DisplayName,
			CreatedAt: administrator.CreatedAt.Format(time.RFC3339), UpdatedAt: administrator.UpdatedAt.Format(time.RFC3339),
		},
		AuthenticatedAt: authenticatedAt.Format(time.RFC3339), IdleExpiresAt: idleExpiresAt.Format(time.RFC3339),
		AbsoluteExpiresAt: absoluteExpiresAt.Format(time.RFC3339),
	}}
}
