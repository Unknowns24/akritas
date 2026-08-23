package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

// CreateRemediationBranchCommand carries a caller-supplied RemediationID so
// branch naming and idempotent replay never depend on IssueReference or a
// GitHub Issue existing.
type CreateRemediationBranchCommand struct {
	RemediationID   uuid.UUID
	IncidentID      uuid.UUID
	InvestigationID uuid.UUID
	WorkspacePath   string
	BaseBranch      string
}

type ExecuteRemediationValidationsCommand struct {
	RemediationID uuid.UUID
	WorkspacePath string
}

type CreateRemediationPullRequestCommand struct {
	RemediationID  uuid.UUID
	WorkspacePath  string
	IdempotencyKey uuid.UUID
	WorkspaceRoot  string
}

type StartIncidentRemediationCommand struct {
	IncidentID     uuid.UUID
	IdempotencyKey uuid.UUID
	WorkspaceRoot  string
}

// RemediationUseCase intentionally excludes automatic change generation,
// merge, deploy and rollback. Pull-request creation is an explicit final
// step and stops after the PR reference is persisted.
type RemediationUseCase interface {
	StartIncidentRemediation(context.Context, StartIncidentRemediationCommand) (*domain.Operation, error)
	GetIncidentRemediation(context.Context, uuid.UUID) (*domain.Remediation, error)
	GetRemediation(context.Context, uuid.UUID) (*domain.Remediation, error)
	ListValidationResults(context.Context, uuid.UUID, paging.Params) (paging.Slice[domain.ValidationResult], error)
	CreateRemediationBranch(context.Context, CreateRemediationBranchCommand) (*domain.Remediation, error)
	ExecuteRemediationValidations(context.Context, ExecuteRemediationValidationsCommand) (*domain.Remediation, []domain.ValidationResult, error)
	QueueRemediationPullRequest(context.Context, CreateRemediationPullRequestCommand) (*domain.Operation, error)
	CreateRemediationPullRequest(context.Context, CreateRemediationPullRequestCommand) (*domain.Remediation, error)
	ListPullRequests(context.Context, paging.Params) (paging.Slice[domain.PullRequestProjection], error)
	GetPullRequest(context.Context, uuid.UUID) (*domain.PullRequestProjection, error)
}
