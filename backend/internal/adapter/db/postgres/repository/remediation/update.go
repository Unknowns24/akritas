package remediation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/txcontext"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (r *Repository) Update(ctx context.Context, value *domain.Remediation) error {
	if value == nil || value.Validate() != nil {
		return domain.ErrInvalidRemediation
	}
	record := fromDomain(value)
	result := txcontext.From(ctx, r.db).WithContext(ctx).Table("remediations").
		Where("id = ?", value.ID).
		Updates(map[string]any{
			"investigation_id":        record.InvestigationID,
			"status":                  record.Status,
			"branch_name":             record.BranchName,
			"changes_summary":         record.ChangesSummary,
			"failure_user_message":    record.FailureUserMessage,
			"pull_request_number":     record.PullRequestNumber,
			"pull_request_url":        record.PullRequestURL,
			"pull_request_repository": record.PullRequestRepository,
			"pull_request_branch":     record.PullRequestBranch,
			"pull_request_created_at": record.PullRequestCreatedAt,
			"updated_at":              record.UpdatedAt,
		})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRemediationNotFound
	}
	return nil
}
