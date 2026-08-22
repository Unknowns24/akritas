package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"gorm.io/gorm"
)

func (r *Repository) CreateWithCredential(ctx context.Context, server *domain.DokployServer, secret portsout.SecretValue) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("dokploy_servers").Create(server).Error; err != nil {
			return mapError(err)
		}
		return r.credentials.PutTx(ctx, tx, portsout.CredentialOwnerDokployServer, server.ID, secret)
	})
}
