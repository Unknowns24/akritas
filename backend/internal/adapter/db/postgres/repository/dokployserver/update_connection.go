package dokployserver

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) UpdateConnection(ctx context.Context, server *domain.DokployServer) error {
	result := r.db.WithContext(ctx).Table("dokploy_servers").Model(&domain.DokployServer{}).Where("id = ?", server.ID).Updates(map[string]any{
		"connection_status": server.ConnectionStatus,
		"application_count": server.ApplicationCount,
		"last_synced_at":    server.LastSyncedAt,
		"updated_at":        server.UpdatedAt,
	})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrIntegrationNotFound
	}
	return nil
}
