package githubapp

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (r *Repository) CreateRegistration(ctx context.Context, registration *portsout.GitHubAppRegistration) error {
	if registration == nil || registration.ID.String() == "" || len(registration.ConversionStateDigest) != 32 || registration.ExpiresAt.After(registration.CreatedAt.Add(time.Hour)) {
		return domain.ErrManifestStateInvalid
	}
	if err := r.db.WithContext(ctx).Table("github_app_registrations").Create(registration).Error; err != nil {
		return mapError(err)
	}
	return nil
}
