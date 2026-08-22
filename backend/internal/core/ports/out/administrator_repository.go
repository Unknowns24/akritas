package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// AdministratorAuthentication carries persistence-only password material
// alongside the safe domain entity. TOTP seeds remain exclusively in the
// CredentialStore.
type AdministratorAuthentication struct {
	Administrator domain.Administrator
	PasswordHash  string
}

type AdministratorRepository interface {
	ExistsActive(ctx context.Context) (bool, error)
	// Create persists the single Administrator. passwordHash is persistence-only
	// material that never belongs on domain.Administrator itself. Implementations must map a
	// unique-email constraint violation to domain.ErrAdministratorAlreadyExists.
	Create(ctx context.Context, administrator *domain.Administrator, passwordHash string) error
	// FindByID returns (nil, nil) when no administrator with that id exists.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Administrator, error)
	// FindByEmail returns (nil, nil) when no administrator with that email
	// exists.
	FindByEmail(ctx context.Context, email string) (*AdministratorAuthentication, error)
	ConsumeTOTPPeriod(ctx context.Context, id uuid.UUID, period int64) (bool, error)
}
