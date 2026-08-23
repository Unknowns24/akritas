package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// CreateRemediationBranchCommand carries a caller-supplied RemediationID so
// branch naming and idempotent replay never depend on IssueReference or a
// GitHub Issue existing.
type CreateRemediationBranchCommand struct {
	RemediationID uuid.UUID
	IncidentID    uuid.UUID
	WorkspacePath string
	BaseBranch    string
}

type ExecuteRemediationValidationsCommand struct {
	RemediationID uuid.UUID
	WorkspacePath string
}

type CreateRemediationPullRequestCommand struct {
	RemediationID uuid.UUID
	WorkspacePath string
}

// RemediationUseCase intentionally excludes automatic change generation,
// merge, deploy and rollback. Pull-request creation is an explicit final
// step and stops after the PR reference is persisted.
type RemediationUseCase interface {
	CreateRemediationBranch(context.Context, CreateRemediationBranchCommand) (*domain.Remediation, error)
	ExecuteRemediationValidations(context.Context, ExecuteRemediationValidationsCommand) (*domain.Remediation, []domain.ValidationResult, error)
	CreateRemediationPullRequest(context.Context, CreateRemediationPullRequestCommand) (*domain.Remediation, error)
}
