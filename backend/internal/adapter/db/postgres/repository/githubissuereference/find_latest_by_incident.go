package githubissuereference

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindLatestByIncident(ctx context.Context, incidentID uuid.UUID) (*domain.GitHubIssueReference, error) {
	var value domain.GitHubIssueReference
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("github_issue_references").
		Where("incident_id = ?", incidentID).Order("created_at DESC, investigation_id DESC").Take(&value).Error
	if err != nil {
		return nil, mapError(err)
	}
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
