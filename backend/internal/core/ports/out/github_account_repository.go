package out

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type GitHubAccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.GitHubAccount, error)
}
