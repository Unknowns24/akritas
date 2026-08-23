package githubissuereference

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Create(ctx context.Context, value *domain.GitHubIssueReference) error {
	if value == nil || value.Validate() != nil {
		return domain.ErrInvalidGitHubIssueReference
	}
	if err := txcontext.From(ctx, r.db).WithContext(ctx).Table("github_issue_references").Create(value).Error; err != nil {
		return mapError(err)
	}
	return nil
}
