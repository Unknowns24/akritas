package project

import (
	"context"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByDokployApplication(ctx context.Context, serverID uuid.UUID, identifier string) (*domain.Project, error) {
	var value domain.Project
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("projects").Where("dokploy_server_id = ? AND application_identifier = ?", serverID, strings.TrimSpace(identifier)).First(&value).Error
	if err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
