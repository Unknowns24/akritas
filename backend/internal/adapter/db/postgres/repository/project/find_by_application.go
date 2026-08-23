package project

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) FindByDokploySource(ctx context.Context, selector domain.DokploySourceSelector) (*domain.Project, error) {
	selector = selector.Normalize()
	var value domain.Project
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where(
		"dokploy_server_id = ? AND source_type = ? AND source_resource_identifier = ? AND COALESCE(source_service_name, '') = ?",
		selector.DokployServerID, selector.Type, selector.ResourceIdentifier, selector.ServiceName,
	).First(&value).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
