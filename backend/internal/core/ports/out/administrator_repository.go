package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// AdministratorCredentials carries persistence-only secret material
// alongside the safe domain.Administrator, the same pattern Create already
// uses: the domain entity itself never carries a password hash, an
// encrypted TOTP secret, or the last accepted TOTP period.
type AdministratorCredentials struct {
	Administrator          domain.Administrator
	PasswordHash           string
	EncryptedTOTPSecret    []byte
	LastAcceptedTOTPPeriod int64
}

type AdministratorRepository interface {
	ExistsActive(ctx context.Context) (bool, error)
	// Create persists the single Administrator. passwordHash and
	// encryptedTOTPSecret are persistence-only secret material that never
	// belongs on domain.Administrator itself. Implementations must map a
	// unique-email constraint violation to domain.ErrAdministratorAlreadyExists.
	Create(ctx context.Context, administrator *domain.Administrator, passwordHash string, encryptedTOTPSecret []byte) error
	// FindByID returns (nil, nil) when no administrator with that id exists.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Administrator, error)
	// FindByEmail returns (nil, nil) when no administrator with that email
	// exists.
	FindByEmail(ctx context.Context, email string) (*AdministratorCredentials, error)
	UpdateLastAcceptedTOTPPeriod(ctx context.Context, id uuid.UUID, period int64) error
}
