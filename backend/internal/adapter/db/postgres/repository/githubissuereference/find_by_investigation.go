package githubissuereference

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) FindByInvestigation(ctx context.Context, investigationID uuid.UUID) (*domain.GitHubIssueReference, error) {
	var value domain.GitHubIssueReference
	err := txcontext.From(ctx, r.db).WithContext(ctx).Table("github_issue_references").
		Where("investigation_id = ?", investigationID).Take(&value).Error
	if err != nil {
		return nil, mapError(err)
	}
	if err := value.Validate(); err != nil {
		return nil, mapError(err)
	}
	return &value, nil
}
