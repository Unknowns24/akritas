package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type AdministratorRepository interface {
	ExistsActive(ctx context.Context) (bool, error)
	// Create persists the single Administrator. passwordHash and
	// encryptedTOTPSecret are persistence-only secret material that never
	// belongs on domain.Administrator itself. Implementations must map a
	// unique-email constraint violation to domain.ErrAdministratorAlreadyExists.
	Create(ctx context.Context, administrator *domain.Administrator, passwordHash string, encryptedTOTPSecret []byte) error
}
