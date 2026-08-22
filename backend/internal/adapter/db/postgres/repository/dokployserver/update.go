package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"gorm.io/gorm"
)

func (r *Repository) Update(ctx context.Context, server *domain.DokployServer, secret *portsout.SecretValue) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table("dokploy_servers").Model(&domain.DokployServer{}).Where("id = ?", server.ID).Select("name", "base_url", "server_identifier", "connection_status", "credential_configured", "application_count", "last_synced_at", "updated_at").Updates(server)
		if result.Error != nil {
			return mapError(result.Error)
		}
		if result.RowsAffected == 0 {
			return domain.ErrIntegrationNotFound
		}
		if secret != nil {
			return r.credentials.PutTx(ctx, tx, portsout.CredentialOwnerDokployServer, server.ID, *secret)
		}
		return nil
	})
}
