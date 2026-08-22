package githubapp

import (
	"context"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) ConsumeConversionState(ctx context.Context, digest []byte, now time.Time) (*portsout.GitHubAppRegistration, error) {
	var output *portsout.GitHubAppRegistration
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record domain.GitHubAppRegistration
		if err := tx.Table("github_app_registrations").Clauses(clause.Locking{Strength: "UPDATE"}).Where("conversion_state_digest = ?", digest).First(&record).Error; err != nil {
			return mapError(err)
		}
		if record.Status != portsout.GitHubAppRegistrationCreated || record.ConversionConsumedAt != nil || !now.Before(record.ExpiresAt) {
			return domain.ErrManifestStateConflict
		}
		record.ConversionConsumedAt = &now
		record.UpdatedAt = now
		if err := tx.Table("github_app_registrations").Save(&record).Error; err != nil {
			return mapError(err)
		}
		output = &record
		return nil
	})
	return output, err
}
